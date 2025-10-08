package inf

import (
	"backend/generation/config"
	"backend/generation/core"
	"backend/generation/generator/shared"
	"backend/generation/helper"
	"errors"
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
)

var AllVars = []string{"p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

// extractPremisesAndConclusions 从模板对中提取共享前提、有效结论和无效结论。
func extractPremisesAndConclusions(template config.InferenceTemplatePair) ([]*core.Node, []*core.Node, []*core.Node, error) {
	// 解析共享前提
	premises := make([]*core.Node, 0)
	for _, spec := range template.Premises {
		node, err := exprPatternSpecToAST(&spec)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to parse premise: %v", err)
		}
		premises = append(premises, node)
	}
	if len(template.Valid) == 0 {
		return nil, nil, nil, errors.New("template has no valid conclusions")
	}
	// 解析有效结论
	validConclusions := make([]*core.Node, 0, len(template.Valid))
	for _, body := range template.Valid {
		node, err := exprPatternSpecToAST(&body.Conclusion)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to parse valid conclusion: %v", err)
		}
		validConclusions = append(validConclusions, node)
	}
	// 解析无效结论
	invalidConclusions := make([]*core.Node, 0)
	for _, invalidBody := range template.Invalid {
		node, err := exprPatternSpecToAST(&invalidBody.Conclusion)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to parse invalid conclusion: %v", err)
		}
		invalidConclusions = append(invalidConclusions, node)
	}

	return premises, validConclusions, invalidConclusions, nil
}

// exprPatternSpecToAST 将单个表达式模式规格解析为 AST 节点。
func exprPatternSpecToAST(spec *config.ExprPatternSpec) (*core.Node, error) {
	if spec == nil {
		return nil, nil
	}

	kindName := strings.ToUpper(spec.Kind)
	kind, ok := config.OpNameToKind[kindName]
	if !ok {
		return nil, errors.New("unknown operator: " + spec.Kind)
	}

	node := &core.Node{Kind: kind}

	switch kind {
	case core.Var:
		node.Name = spec.Name
		return node, nil
	case core.Not:
		left, err := exprPatternSpecToAST(spec.Left)
		if err != nil {
			return nil, err
		}
		node.Left = left
		return node, nil
	default:
		left, err := exprPatternSpecToAST(spec.Left)
		if err != nil {
			return nil, err
		}
		right, err := exprPatternSpecToAST(spec.Right)
		if err != nil {
			return nil, err
		}
		node.Left = left
		node.Right = right
		return node, nil
	}
}

// instantiateExprPattern 用绑定替换表达式模式中的占位符，生成具体的 AST。
func instantiateExprPattern(pattern *core.Node, binding map[string]*core.Node) *core.Node {
	if pattern == nil {
		return nil
	}
	if pattern.Kind == core.Var {
		if node, ok := binding[pattern.Name]; ok {
			return node.Clone()
		}
		// 未绑定占位符的情况，根据需要返回 nil/报错
		return nil
	}
	left := instantiateExprPattern(pattern.Left, binding)
	right := instantiateExprPattern(pattern.Right, binding)
	return &core.Node{Kind: pattern.Kind, Left: left, Right: right}
}

// prepareBindings 根据模板的槽位规格生成槽位绑定。
// 推理模板的变量需求完全由模板自身决定，这里直接按字母顺序分配变量，
// 不受 profile.Vars 限制
func (g InferenceGenerator) prepareBindings(rng *rand.Rand, slots map[string]config.InferenceTemplateSlot) (map[string]*core.Node, error) {
	// 按名字排序槽位，确保绑定顺序一致
	keys := make([]string, 0, len(slots))
	for name := range slots {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	idx := 0
	nextVar := func() string {
		name := AllVars[idx]
		idx++
		return name
	}
	buildSlot := func(strategy string) *core.Node {
		switch strategy {
		case "use_var":
			return shared.NewVar(nextVar())
		case "use_neg_var":
			return shared.Unary(core.Not, shared.NewVar(nextVar()))
		case "use_and":
			left := nextVar()
			right := nextVar()
			return shared.Binary(core.And, shared.NewVar(left), shared.NewVar(right))
		case "use_or":
			left := nextVar()
			right := nextVar()
			return shared.Binary(core.Or, shared.NewVar(left), shared.NewVar(right))
		case "use_const_true":
			return shared.TautologyFromVar(nextVar())
		case "use_const_false":
			return shared.ContradictionFromVar(nextVar())
		default:
			return shared.NewVar(nextVar())
		}
	}
	bindings := make(map[string]*core.Node, len(slots))

	for _, name := range keys {
		slot := slots[name]
		fillerDist, ok := g.cfg.SlotFillers[slot.Type]
		if !ok {
			return nil, fmt.Errorf("prepareBindings: slot type %q missing slot filler spec", slot.Type)
		}
		// 按权重采样, 默认使用变量
		choice := helper.SampleWeighted(fillerDist, rng)
		if choice == "" {
			choice = "use_var"
		}
		node := buildSlot(choice)
		bindings[name] = node
	}

	return bindings, nil
}

// instantiateTemplatePair 用绑定替换模板对中的占位符，生成具体的前提和结论。
func instantiateTemplatePair(premises []*core.Node, validCons []*core.Node, inValidCons []*core.Node, bindings map[string]*core.Node) ([]*core.Node, []*core.Node, []*core.Node) {
	// 替换premises 中的占位符
	preInst := make([]*core.Node, 0, len(premises))
	for _, p := range premises {
		inst := instantiateExprPattern(p, bindings)
		preInst = append(preInst, inst)
	}
	// 替换 validCons 中的占位符
	validInst := make([]*core.Node, 0, len(validCons))
	for _, vc := range validCons {
		inst := instantiateExprPattern(vc, bindings)
		validInst = append(validInst, inst)
	}
	// 替换 inValidCons 中的占位符
	inValidInst := make([]*core.Node, 0, len(inValidCons))
	for _, ic := range inValidCons {
		inst := instantiateExprPattern(ic, bindings)
		inValidInst = append(inValidInst, inst)
	}
	return preInst, validInst, inValidInst

}

// stringifyInstances 将前提和结论实例化后的 AST 转换为字符串表示。
func stringifyInstances(premises []*core.Node, validCons []*core.Node, inValidCons []*core.Node) (string, []string, []string) {
	premiseStrs := make([]string, 0, len(premises))
	for _, p := range premises {
		premiseStrs = append(premiseStrs, helper.Stringify(p))
	}
	validStrs := make([]string, 0, len(validCons))
	for _, vc := range validCons {
		validStrs = append(validStrs, helper.Stringify(vc))
	}
	inValidStrs := make([]string, 0, len(inValidCons))
	for _, ic := range inValidCons {
		inValidStrs = append(inValidStrs, helper.Stringify(ic))
	}
	return strings.Join(premiseStrs, ", "), validStrs, inValidStrs
}

// 统计 AST 中的变量
func collectVars(nodes []*core.Node) []string {
	varSet := make(map[string]struct{})
	// 遍历所有节点，收集变量
	for _, node := range nodes {
		vars := shared.FilterVars(AllVars, node)
		for _, v := range vars {
			// 记录变量
			varSet[v] = struct{}{}
		}
	}
	vars := make([]string, 0, len(varSet))
	for v := range varSet {
		vars = append(vars, v)
	}
	return vars
}

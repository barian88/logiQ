package inf

import (
	"backend/generation/core"
	"backend/generation/generator/eq"
	"backend/generation/generator/shared"
	"backend/generation/helper"
	"math/rand/v2"
)

// 每个原始node有一定的概率变形，生成一个包含自身的候选池，池子通过对结论本身进行可用等价变换得到，然后随机选择一个作为node。

// 建立等价规则表
var (
	ruleGroup = []eq.Rule{
		{Name: "demorgan", Apply: eq.WrapDeterministic(eq.DeMorganVariants)},
		{Name: "distribute", Apply: eq.WrapDeterministic(eq.DistributivityVariants)},
		{Name: "implication", Apply: eq.WrapDeterministic(eq.ImplicationVariants)},
		{Name: "biconditional", Apply: eq.WrapDeterministic(eq.BiconditionalVariants)},
	}
)

// transformNodes 对node集合进行变形
// oriNodes: 原始node集合
// rng: 随机数生成器
// 返回变形后的集合
func transformNodes(oriNodes []*core.Node, rng *rand.Rand) []*core.Node {
	// 复用eq包中的RuleExecutor
	// 选择变形规则
	limit := 3 // 每个node最多变形3次
	rules := ruleGroup
	ruleExecutor := eq.RuleExecutor{Rules: rules, Limit: limit, Rng: rng}
	// 建立一个变形后的集合
	transformedNodes := make([]*core.Node, 0, len(oriNodes))
	// 对每个node进行变形
	for _, node := range oriNodes {
		if rng.Float64() < 0.75 {
			variants := ruleExecutor.Execute(node)
			// 把自身也加入候选池
			variants = append(variants, node.Clone())
			// 从候选池中随机选择一个
			chosen := variants[rng.IntN(len(variants))]
			transformedNodes = append(transformedNodes, chosen)
		} else {
			// 保持不变
			transformedNodes = append(transformedNodes, node.Clone())
		}
	}
	// 返回变形后的集合
	return transformedNodes
}

// expandConclusions 对有效结论和无效结论进行进一步可能的扩展
// 原理是，如果A and B 在有效结论中，那么A和B单独出现也应该是有效结论
// 如果A or B 在无效结论中，那么A和B单独出现也应该是无效结论
// 其他情况不做处理
func expandConclusions(valid, invalid []*core.Node) ([]*core.Node, []*core.Node) {
	expandedValid := make([]*core.Node, 0)
	for _, node := range valid {
		expandedValid = append(expandedValid, node.Clone())
		if node.Kind == core.And && node.Left != nil && node.Right != nil {
			expandedValid = append(expandedValid, node.Left.Clone(), node.Right.Clone())
		}
	}

	expandedInvalid := make([]*core.Node, 0)
	for _, node := range invalid {
		expandedInvalid = append(expandedInvalid, node.Clone())
		if node.Kind == core.Or && node.Left != nil && node.Right != nil {
			expandedInvalid = append(expandedInvalid, node.Left.Clone(), node.Right.Clone())
		}
	}
	return expandedValid, expandedInvalid
}

// 因为变形后可能会导致结论和前提重复，所以需要过滤掉前提中已经出现的结论
func filterContains(premises []*core.Node, target []*core.Node) []*core.Node {
	filtered := make([]*core.Node, 0, len(target))
	premiseSet := make(map[string]struct{})
	for _, p := range premises {
		premiseSet[helper.Stringify(p)] = struct{}{}
	}
	for _, t := range target {
		if _, exists := premiseSet[helper.Stringify(t)]; !exists {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

// ExpandAndTransform 对前提、有效结论和无效结论进行变形和扩展.
// 一个统一的入口，封装了上面的三个方法
func ExpandAndTransform(premises, validConclusions, invalidConclusions []*core.Node, rng *rand.Rand) ([]*core.Node, []*core.Node, []*core.Node) {
	//先变形
	transformedPremises := transformNodes(premises, rng)
	transformedValid := transformNodes(validConclusions, rng)
	transformedInvalid := transformNodes(invalidConclusions, rng)
	//再扩展
	expandedValid, expandedInvalid := expandConclusions(transformedValid, transformedInvalid)
	//再过滤
	filterValid := filterContains(transformedPremises, expandedValid)
	filterInvalid := filterContains(transformedPremises, expandedInvalid)
	// 去重
	finalValid := shared.DedupNodes(filterValid)
	finalInvalid := shared.DedupNodes(filterInvalid)

	return transformedPremises, finalValid, finalInvalid
}

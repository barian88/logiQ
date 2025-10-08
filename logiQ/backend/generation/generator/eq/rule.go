package eq

import (
	"backend/generation/config"
	"backend/generation/core"
	"math/rand/v2"
)

// Rule 表示一条可用于生成公式变体的变换规则。
type Rule struct {
	Name  string
	Apply func(*core.Node, *rand.Rand) []*core.Node
}

type ruleGroup int

const (
	ruleGroupEquivalent ruleGroup = iota
	ruleGroupNonEquivalent
)

// GenerateEquivalentVariants 根据难度配置生成与根公式等价的候选。
func GenerateEquivalentVariants(root *core.Node, rng *rand.Rand, cfg config.EquivalenceConfig, chainStep int) []*core.Node {
	return generateVariants(root, rng, cfg, chainStep, ruleGroupEquivalent)
}

// GenerateNonEquivalentVariants 根据难度配置生成与根公式不等价的候选。
func GenerateNonEquivalentVariants(root *core.Node, rng *rand.Rand, cfg config.EquivalenceConfig, chainStep int) []*core.Node {
	return generateVariants(root, rng, cfg, chainStep, ruleGroupNonEquivalent)
}

func generateVariants(root *core.Node, rng *rand.Rand, cfg config.EquivalenceConfig, chainStep int, group ruleGroup) []*core.Node {

	rules := collectConfiguredRules(cfg, group)
	executor := RuleExecutor{Rng: rng, Rules: rules, Limit: chainStep}
	return executor.Execute(root)
}

func collectConfiguredRules(cfg config.EquivalenceConfig, group ruleGroup) []Rule {
	registry := builtinRuleRegistry[group]
	if len(registry) == 0 {
		return nil
	}

	specs := group.specs(cfg)
	if len(specs) == 0 {
		return nil
	}

	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		if rule, ok := registry[spec.Name]; ok {
			rules = append(rules, rule)
		}
	}
	return rules
}

func (g ruleGroup) specs(cfg config.EquivalenceConfig) []config.EquivalenceRuleSpec {
	switch g {
	case ruleGroupEquivalent:
		return cfg.Rules.Equivalent
	case ruleGroupNonEquivalent:
		return cfg.Rules.NonEquivalent
	default:
		return nil
	}
}

var builtinRuleRegistry = map[ruleGroup]map[string]Rule{
	ruleGroupEquivalent: {
		"double_negation":         {Name: "double_negation", Apply: WrapDeterministic(doubleNegationVariants)},
		"implication_elimination": {Name: "implication_elimination", Apply: WrapDeterministic(ImplicationVariants)},
		"commutativity":           {Name: "commutativity", Apply: WrapDeterministic(commutativityVariants)},
		"de_morgan":               {Name: "de_morgan", Apply: WrapDeterministic(DeMorganVariants)},
		"biconditional_expansion": {Name: "biconditional_expansion", Apply: WrapDeterministic(BiconditionalVariants)},
		"associativity":           {Name: "associativity", Apply: WrapDeterministic(associativityVariants)},
		"distributivity":          {Name: "distributivity", Apply: WrapDeterministic(DistributivityVariants)},
	},
	ruleGroupNonEquivalent: {
		"negate_root":         {Name: "negate_root", Apply: WrapDeterministic(negateRootVariants)},
		"reverse_implication": {Name: "reverse_implication", Apply: WrapDeterministic(reverseImplicationVariants)},
		"flip_operator":       {Name: "flip_operator", Apply: WrapDeterministic(flipOperatorVariants)},
		"mutate_literal":      {Name: "mutate_literal", Apply: wrapWithRNG(mutateLiteralVariants)},
	},
}

func WrapDeterministic(fn func(*core.Node) []*core.Node) func(*core.Node, *rand.Rand) []*core.Node {
	return func(node *core.Node, _ *rand.Rand) []*core.Node {
		return fn(node)
	}
}

func wrapWithRNG(fn func(*core.Node, *rand.Rand) []*core.Node) func(*core.Node, *rand.Rand) []*core.Node {
	return fn
}

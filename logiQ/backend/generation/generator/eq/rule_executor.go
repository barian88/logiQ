package eq

import (
	"backend/generation/core"
	"backend/generation/generator/shared"
	"backend/generation/helper"
	"math/rand/v2"
)

type RuleExecutor struct {
	Rng   *rand.Rand
	Rules []Rule
	Limit int
}

type ruleState struct {
	node *core.Node
	used stringSet
}

type stringSet map[string]struct{}

func (e RuleExecutor) Execute(root *core.Node) []*core.Node {
	if root == nil || e.Limit <= 0 || len(e.Rules) == 0 {
		return nil
	}

	// results 保存已生成的唯一公式候选（按字符串签名去重）
	results := make(map[string]*core.Node)
	// frontier 维护当前深度可继续扩展的状态集合
	frontier := []ruleState{{node: root.Clone(), used: stringSet{}}}

	for step := 0; step < e.Limit; step++ {
		// next 收集下一层待扩展状态，seen 用于避免同层重复
		next := make([]ruleState, 0)
		// 记录本层已见公式，避免重复收集.比遍历 next 快
		seen := stringSet{}

		for _, st := range frontier {
			for _, rule := range e.Rules {
				if st.used.contains(rule.Name) {
					continue
				}

				for _, variant := range applyRuleRecursive(st.node, e.Rng, rule.Apply) {
					if variant == nil {
						continue
					}

					// 记录新公式，避免重复收集
					sig := helper.Stringify(variant)
					if _, ok := results[sig]; !ok {
						results[sig] = variant
					}
					// 达到步数上限或本层已见则跳过
					if step+1 >= e.Limit {
						continue
					}
					if seen.contains(sig) {
						continue
					}

					// 将当前规则标记入 used，供后续层判断是否复用
					seen.add(sig)
					next = append(next, ruleState{node: variant, used: st.used.cloneWith(rule.Name)})
				}
			}
		}

		if len(next) == 0 {
			break
		}
		frontier = next
	}

	// 生成最终列表并再做一次去重
	out := make([]*core.Node, 0, len(results))
	for _, node := range results {
		out = append(out, node)
	}
	return shared.DedupNodes(out)
}

func (s stringSet) contains(val string) bool {
	if s == nil {
		return false
	}
	_, ok := s[val]
	return ok
}

func (s stringSet) add(val string) {
	if s == nil {
		return
	}
	s[val] = struct{}{}
}

func (s stringSet) cloneWith(val string) stringSet {
	out := make(stringSet, len(s)+1)
	for existing := range s {
		out[existing] = struct{}{}
	}
	out[val] = struct{}{}
	return out
}

func applyRuleRecursive(node *core.Node, rng *rand.Rand, apply func(*core.Node, *rand.Rand) []*core.Node) []*core.Node {
	if node == nil {
		return nil
	}
	var variants []*core.Node
	for _, replacement := range apply(node, rng) {
		variants = append(variants, replacement)
	}
	if node.Left != nil {
		for _, leftVariant := range applyRuleRecursive(node.Left, rng, apply) {
			clone := node.Clone()
			clone.Left = leftVariant
			variants = append(variants, clone)
		}
	}
	if node.Right != nil {
		for _, rightVariant := range applyRuleRecursive(node.Right, rng, apply) {
			clone := node.Clone()
			clone.Right = rightVariant
			variants = append(variants, clone)
		}
	}
	return variants
}

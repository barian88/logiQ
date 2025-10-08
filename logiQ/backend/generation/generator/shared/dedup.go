package shared

import (
	"backend/generation/core"
	"backend/generation/helper"
)

func DedupNodes(nodes []*core.Node) []*core.Node {
	if len(nodes) <= 1 {
		return nodes
	}
	seen := make(map[string]struct{}, len(nodes))
	result := make([]*core.Node, 0, len(nodes))
	for _, n := range nodes {
		if n == nil {
			continue
		}
		sig := helper.Stringify(n)
		if _, ok := seen[sig]; ok {
			continue
		}
		seen[sig] = struct{}{}
		result = append(result, n)
	}
	return result
}

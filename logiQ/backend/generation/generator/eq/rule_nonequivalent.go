package eq

import (
	"backend/generation/core"
	"math/rand/v2"
)

func reverseImplicationVariants(node *core.Node) []*core.Node {
	if node.Kind != core.Impl || node.Left == nil || node.Right == nil {
		return nil
	}
	return []*core.Node{{Kind: core.Impl, Left: node.Right.Clone(), Right: node.Left.Clone()}}
}

func negateRootVariants(node *core.Node) []*core.Node {
	return []*core.Node{{Kind: core.Not, Left: node.Clone()}}
}

func flipOperatorVariants(node *core.Node) []*core.Node {
	if node.Left == nil || node.Right == nil {
		return nil
	}
	switch node.Kind {
	case core.And:
		return []*core.Node{{Kind: core.Or, Left: node.Left.Clone(), Right: node.Right.Clone()}}
	case core.Or:
		return []*core.Node{{Kind: core.And, Left: node.Left.Clone(), Right: node.Right.Clone()}}
	case core.Impl:
		return []*core.Node{{Kind: core.And, Left: node.Left.Clone(), Right: node.Right.Clone()}}
	case core.Iff:
		return []*core.Node{{Kind: core.Impl, Left: node.Left.Clone(), Right: node.Right.Clone()}}
	default:
		return nil
	}
}

func mutateLiteralVariants(node *core.Node, rng *rand.Rand) []*core.Node {
	switch node.Kind {
	case core.Var:
		return []*core.Node{{Kind: core.Not, Left: node.Clone()}}
	case core.Not:
		if node.Left != nil && node.Left.Kind == core.Var {
			return []*core.Node{node.Left.Clone()}
		}
	}
	if node.Left == nil && node.Right == nil {
		return nil
	}
	if node.Left != nil && node.Right != nil {
		clone := node.Clone()
		if rng != nil && rng.IntN(2) == 1 {
			right := clone.Right
			clone.Right = &core.Node{Kind: core.Not, Left: right}
		} else {
			left := clone.Left
			clone.Left = &core.Node{Kind: core.Not, Left: left}
		}
		return []*core.Node{clone}
	}
	if node.Left != nil {
		clone := node.Clone()
		left := clone.Left
		clone.Left = &core.Node{Kind: core.Not, Left: left}
		return []*core.Node{clone}
	}
	if node.Right != nil {
		clone := node.Clone()
		right := clone.Right
		clone.Right = &core.Node{Kind: core.Not, Left: right}
		return []*core.Node{clone}
	}
	return nil
}

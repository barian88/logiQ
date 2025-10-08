package eq

import (
	"backend/generation/core"
	"backend/generation/generator/shared"
	"backend/generation/helper"
)

func doubleNegationVariants(node *core.Node) []*core.Node {
	var out []*core.Node
	if node.Kind == core.Not && node.Left != nil && node.Left.Kind == core.Not && node.Left.Left != nil {
		out = append(out, node.Left.Left.Clone())
	}
	out = append(out, &core.Node{Kind: core.Not, Left: &core.Node{Kind: core.Not, Left: node.Clone()}})
	return shared.DedupNodes(out)
}

func commutativityVariants(node *core.Node) []*core.Node {
	if node.Left == nil || node.Right == nil {
		return nil
	}
	if node.Kind != core.And && node.Kind != core.Or && node.Kind != core.Iff {
		return nil
	}
	clone := node.Clone()
	clone.Left, clone.Right = clone.Right, clone.Left
	return []*core.Node{clone}
}

func DeMorganVariants(node *core.Node) []*core.Node {
	if node.Kind != core.Not || node.Left == nil {
		return nil
	}
	inner := node.Left
	switch inner.Kind {
	case core.And:
		return []*core.Node{{
			Kind:  core.Or,
			Left:  &core.Node{Kind: core.Not, Left: inner.Left.Clone()},
			Right: &core.Node{Kind: core.Not, Left: inner.Right.Clone()},
		}}
	case core.Or:
		return []*core.Node{{
			Kind:  core.And,
			Left:  &core.Node{Kind: core.Not, Left: inner.Left.Clone()},
			Right: &core.Node{Kind: core.Not, Left: inner.Right.Clone()},
		}}
	default:
		return nil
	}
}

func ImplicationVariants(node *core.Node) []*core.Node {
	switch node.Kind {
	case core.Impl:
		return []*core.Node{{
			Kind:  core.Or,
			Left:  &core.Node{Kind: core.Not, Left: node.Left.Clone()},
			Right: node.Right.Clone(),
		}}
	case core.Or:
		if node.Left != nil && node.Left.Kind == core.Not {
			return []*core.Node{{
				Kind:  core.Impl,
				Left:  node.Left.Left.Clone(),
				Right: node.Right.Clone(),
			}}
		}
	}
	return nil
}

func BiconditionalVariants(node *core.Node) []*core.Node {
	if node.Kind == core.Iff {
		leftImp := &core.Node{Kind: core.Impl, Left: node.Left.Clone(), Right: node.Right.Clone()}
		rightImp := &core.Node{Kind: core.Impl, Left: node.Right.Clone(), Right: node.Left.Clone()}
		return []*core.Node{{Kind: core.And, Left: leftImp, Right: rightImp}}
	}
	if node.Kind == core.And && node.Left != nil && node.Right != nil && node.Left.Kind == core.Impl && node.Right.Kind == core.Impl {
		if helper.Stringify(node.Left.Left) == helper.Stringify(node.Right.Right) && helper.Stringify(node.Left.Right) == helper.Stringify(node.Right.Left) {
			return []*core.Node{{Kind: core.Iff, Left: node.Left.Left.Clone(), Right: node.Left.Right.Clone()}}
		}
	}
	return nil
}

func associativityVariants(node *core.Node) []*core.Node {
	if node.Left == nil || node.Right == nil {
		return nil
	}
	var variants []*core.Node
	switch node.Kind {
	case core.And, core.Or:
		if node.Right.Kind == node.Kind && node.Right.Left != nil && node.Right.Right != nil {
			left := &core.Node{Kind: node.Kind, Left: node.Left.Clone(), Right: node.Right.Left.Clone()}
			variant := &core.Node{Kind: node.Kind, Left: left, Right: node.Right.Right.Clone()}
			variants = append(variants, variant)
		}
		if node.Left.Kind == node.Kind && node.Left.Left != nil && node.Left.Right != nil {
			right := &core.Node{Kind: node.Kind, Left: node.Left.Right.Clone(), Right: node.Right.Clone()}
			variant := &core.Node{Kind: node.Kind, Left: node.Left.Left.Clone(), Right: right}
			variants = append(variants, variant)
		}
	}
	return shared.DedupNodes(variants)
}

func DistributivityVariants(node *core.Node) []*core.Node {
	var variants []*core.Node
	switch node.Kind {
	case core.And:
		if node.Left != nil && node.Right != nil {
			if node.Right.Kind == core.Or && node.Right.Left != nil && node.Right.Right != nil {
				left := &core.Node{Kind: core.And, Left: node.Left.Clone(), Right: node.Right.Left.Clone()}
				right := &core.Node{Kind: core.And, Left: node.Left.Clone(), Right: node.Right.Right.Clone()}
				variant := &core.Node{Kind: core.Or, Left: left, Right: right}
				variants = append(variants, variant)
			}
			if node.Left.Kind == core.Or && node.Left.Left != nil && node.Left.Right != nil {
				left := &core.Node{Kind: core.And, Left: node.Right.Clone(), Right: node.Left.Left.Clone()}
				right := &core.Node{Kind: core.And, Left: node.Right.Clone(), Right: node.Left.Right.Clone()}
				variant := &core.Node{Kind: core.Or, Left: left, Right: right}
				variants = append(variants, variant)
			}
		}
	case core.Or:
		if node.Left != nil && node.Right != nil {
			if node.Right.Kind == core.And && node.Right.Left != nil && node.Right.Right != nil {
				left := &core.Node{Kind: core.Or, Left: node.Left.Clone(), Right: node.Right.Left.Clone()}
				right := &core.Node{Kind: core.Or, Left: node.Left.Clone(), Right: node.Right.Right.Clone()}
				variant := &core.Node{Kind: core.And, Left: left, Right: right}
				variants = append(variants, variant)
			}
			if node.Left.Kind == core.And && node.Left.Left != nil && node.Left.Right != nil {
				left := &core.Node{Kind: core.Or, Left: node.Right.Clone(), Right: node.Left.Left.Clone()}
				right := &core.Node{Kind: core.Or, Left: node.Right.Clone(), Right: node.Left.Right.Clone()}
				variant := &core.Node{Kind: core.And, Left: left, Right: right}
				variants = append(variants, variant)
			}
		}
	}
	if node.Kind == core.Or && node.Left != nil && node.Right != nil {
		if node.Left.Kind == core.And && node.Right.Kind == core.And {
			if common, restLeft, restRight, ok := findCommonFactor(node.Left, node.Right); ok {
				inner := &core.Node{Kind: core.Or, Left: restLeft, Right: restRight}
				variant := &core.Node{Kind: core.And, Left: common, Right: inner}
				variants = append(variants, variant)
			}
		}
	}
	if node.Kind == core.And && node.Left != nil && node.Right != nil {
		if node.Left.Kind == core.Or && node.Right.Kind == core.Or {
			if common, restLeft, restRight, ok := findCommonFactor(node.Left, node.Right); ok {
				inner := &core.Node{Kind: core.And, Left: restLeft, Right: restRight}
				variant := &core.Node{Kind: core.Or, Left: common, Right: inner}
				variants = append(variants, variant)
			}
		}
	}
	return shared.DedupNodes(variants)
}

func findCommonFactor(a, b *core.Node) (common *core.Node, restA *core.Node, restB *core.Node, ok bool) {
	if a == nil || b == nil {
		return nil, nil, nil, false
	}
	children := func(node *core.Node) []*core.Node {
		if node == nil {
			return nil
		}
		return []*core.Node{node.Left, node.Right}
	}
	aChildren := children(a)
	bChildren := children(b)
	for _, candidate := range aChildren {
		if candidate == nil {
			continue
		}
		sig := helper.Stringify(candidate)
		for _, other := range bChildren {
			if other == nil {
				continue
			}
			if helper.Stringify(other) == sig {
				restA = otherChild(a, candidate)
				restB = otherChild(b, other)
				if restA == nil || restB == nil {
					continue
				}
				return candidate.Clone(), restA.Clone(), restB.Clone(), true
			}
		}
	}
	return nil, nil, nil, false
}

func otherChild(parent, child *core.Node) *core.Node {
	if parent == nil {
		return nil
	}
	childSig := helper.Stringify(child)
	if parent.Left != nil && helper.Stringify(parent.Left) == childSig {
		return parent.Right
	}
	if parent.Right != nil && helper.Stringify(parent.Right) == childSig {
		return parent.Left
	}
	return nil
}

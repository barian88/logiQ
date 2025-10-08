package helper

import (
	"backend/generation/core"
)

// Stringify renders the AST into a canonical, fully parenthesised string.
func Stringify(node *core.Node) string {
	if node == nil {
		return ""
	}
	return stringify(node, false)
}

func stringify(node *core.Node, wrap bool) string {
	switch node.Kind {
	case core.Var:
		return node.Name
	case core.Not:
		inner := stringify(node.Left, needsParensForNot(node.Left))
		return "¬" + inner
	case core.And, core.Or, core.Impl, core.Iff:
		left := stringify(node.Left, true)
		right := stringify(node.Right, true)
		op := binaryOpSymbol(node.Kind)
		expr := left + " " + op + " " + right
		if wrap {
			return "(" + expr + ")"
		}
		return expr
	default:
		return ""
	}
}

func needsParensForNot(node *core.Node) bool {
	if node == nil {
		return false
	}
	return node.Kind != core.Var && node.Kind != core.Not
}

func binaryOpSymbol(kind core.NodeKind) string {
	switch kind {
	case core.And:
		return "∧"
	case core.Or:
		return "∨"
	case core.Impl:
		return "→"
	case core.Iff:
		return "↔"
	default:
		return "?"
	}
}

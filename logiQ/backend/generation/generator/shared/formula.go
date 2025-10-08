package shared

import (
	"backend/generation/core"
	"backend/generation/sampler"
	"math/rand/v2"
)

// CanonicalVars returns a deterministic prefix of [p,q,r,s,t] sized to `count` (>=1).
func CanonicalVars(count int) []string {
	base := []string{"p", "q", "r", "s", "t"}
	if count <= 0 {
		count = 1
	}
	if count > len(base) {
		count = len(base)
	}
	return append([]string(nil), base[:count]...)
}

// NewVar returns a variable node with the provided symbol.
func NewVar(name string) *core.Node {
	return &core.Node{Kind: core.Var, Name: name}
}

// Unary builds a unary node (currently only NOT is meaningful).
func Unary(kind core.NodeKind, child *core.Node) *core.Node {
	return &core.Node{Kind: kind, Left: child}
}

// Binary builds a binary node (AND / OR / IMP / IFF).
func Binary(kind core.NodeKind, left, right *core.Node) *core.Node {
	return &core.Node{Kind: kind, Left: left, Right: right}
}

// TautologyFromVar returns a formula akin to P ∨ ¬P using the given symbol.
func TautologyFromVar(name string) *core.Node {
	varNode := NewVar(name)
	return Binary(core.Or, varNode, Unary(core.Not, NewVar(name)))
}

// ContradictionFromVar returns a formula akin to P ∧ ¬P using the given symbol.
func ContradictionFromVar(name string) *core.Node {
	varNode := NewVar(name)
	return Binary(core.And, varNode, Unary(core.Not, NewVar(name)))
}

// RandomFormula synthesises a random propositional formula that respects the
// sampler profile (variable budget, maximum depth, allowed operators).
// It retries until the generated tree contains at least the minimum amount of
// distinct variables required by the profile (two when more than one variable
// is available, otherwise one). A best-effort formula is returned once the
// retry budget is exhausted.
func RandomFormula(rng *rand.Rand, prof sampler.Profile) *core.Node {
	vars := CanonicalVars(prof.Vars)
	if len(vars) == 0 {
		vars = []string{"p"}
	}

	allowed := prof.AllowedOps
	if len(allowed) == 0 {
		allowed = []core.NodeKind{core.And, core.Or, core.Not}
	}
	const maxAttempts = 64
	needDistinct := len(vars) >= 2
	minVarCount := 1
	if needDistinct {
		minVarCount = 2
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		root := generateNode(rng, allowed, prof.MaxDepth, vars)
		if root == nil {
			continue
		}
		used := make(map[string]struct{})
		collectVars(root, used)
		if len(used) < minVarCount {
			continue
		}
		return root
	}

	return generateNode(rng, allowed, prof.MaxDepth, vars)
}

// FilterVars returns only the variables from `all` which appear in `formula`, preserving order.
func FilterVars(all []string, formula *core.Node) []string {
	used := make(map[string]struct{})
	collectVars(formula, used)
	if len(used) == 0 {
		return all
	}
	ordered := make([]string, 0, len(used))
	for _, v := range all {
		if _, ok := used[v]; ok {
			ordered = append(ordered, v)
		}
	}
	return ordered
}

// generateNode recursively builds an AST using the provided operator set.
// Leaf nodes are variables; inner nodes are chosen uniformly from the allowed
// operators while respecting the remaining depth budget.
func generateNode(rng *rand.Rand, ops []core.NodeKind, maxDepth int, vars []string) *core.Node {
	if rng == nil {
		return nil
	}
	if maxDepth <= 0 {
		return randomVarNode(rng, vars)
	}
	if maxDepth == 1 || rng.Float64() < 0.1 {
		return randomVarNode(rng, vars)
	}

	op := ops[rng.IntN(len(ops))]
	switch op {
	case core.Not:
		return &core.Node{Kind: core.Not, Left: generateNode(rng, ops, maxDepth-1, vars)}
	case core.And, core.Or, core.Impl, core.Iff:
		left := generateNode(rng, ops, maxDepth-1, vars)
		right := generateNode(rng, ops, maxDepth-1, vars)
		return &core.Node{Kind: op, Left: left, Right: right}
	default:
		return randomVarNode(rng, vars)
	}
}

// randomVarNode draws one variable symbol and yields it as an AST leaf.
func randomVarNode(rng *rand.Rand, vars []string) *core.Node {
	if len(vars) == 0 {
		return &core.Node{Kind: core.Var, Name: "p"}
	}
	name := vars[rng.IntN(len(vars))]
	return &core.Node{Kind: core.Var, Name: name}
}

// collectVars walks the AST and records every variable symbol it encounters.
func collectVars(node *core.Node, set map[string]struct{}) {
	if node == nil {
		return
	}
	if node.Kind == core.Var {
		set[node.Name] = struct{}{}
	}
	collectVars(node.Left, set)
	collectVars(node.Right, set)
}

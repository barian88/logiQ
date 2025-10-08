package inf

import "backend/generation/core"

type ExprPattern struct {
	Kind core.NodeKind
	Vars []string
}

type Template struct {
	Name       string
	ChainSteps int
	Premises   []ExprPattern
	Conclusion ExprPattern
}

var validTemplates = []Template{
	{
		Name:       "Modus Ponens",
		ChainSteps: 1,
		Premises: []ExprPattern{
			{Kind: core.Impl, Vars: []string{"A", "B"}},
			{Kind: core.Var, Vars: []string{"A"}},
		},
		Conclusion: ExprPattern{Kind: core.Var, Vars: []string{"B"}},
	},
}

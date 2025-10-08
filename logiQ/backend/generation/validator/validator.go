package validator

import (
	"backend/generation/core"
	"backend/generation/helper"
)

// Validator 这个接口目前没有用到，defaultValidator满足了需求
type Validator interface {
	Eval(formula *core.Node, assign map[string]bool) bool
	Equivalent(a, b *core.Node, vars []string) bool
	Derivable(premises []*core.Node, concl *core.Node, vars []string) bool
}

type DefaultValidator struct{}

func NewDefaultValidator() DefaultValidator { return DefaultValidator{} }

func (v DefaultValidator) Eval(formula *core.Node, assign map[string]bool) bool {
	if formula == nil {
		return false
	}
	switch formula.Kind {
	case core.Var:
		return lookup(assign, formula.Name)
	case core.Not:
		return !v.Eval(formula.Left, assign)
	case core.And:
		return v.Eval(formula.Left, assign) && v.Eval(formula.Right, assign)
	case core.Or:
		return v.Eval(formula.Left, assign) || v.Eval(formula.Right, assign)
	case core.Impl:
		return !v.Eval(formula.Left, assign) || v.Eval(formula.Right, assign)
	case core.Iff:
		left := v.Eval(formula.Left, assign)
		right := v.Eval(formula.Right, assign)
		return left == right
	default:
		return false
	}
}

func (v DefaultValidator) Equivalent(a, b *core.Node, vars []string) bool {
	assignments := helper.EnumerateAssignments(vars)
	for _, assign := range assignments {
		if v.Eval(a, assign) != v.Eval(b, assign) {
			return false
		}
	}
	return true
}

func (v DefaultValidator) Derivable(premises []*core.Node, concl *core.Node, vars []string) bool {
	hasSupportingAssignment := false
	assignments := helper.EnumerateAssignments(vars)

	for _, assign := range assignments {
		allPremisesTrue := true
		for _, prem := range premises {
			if !v.Eval(prem, assign) {
				allPremisesTrue = false
				break
			}
		}

		if !allPremisesTrue {
			continue
		}

		hasSupportingAssignment = true
		if !v.Eval(concl, assign) {
			return false
		}
	}

	return hasSupportingAssignment
}

func lookup(assign map[string]bool, name string) bool {
	if assign == nil {
		return false
	}
	val, ok := assign[name]
	if !ok {
		return false
	}
	return val
}

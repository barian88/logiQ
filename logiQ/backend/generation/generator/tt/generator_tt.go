package tt

import (
	"backend/generation/core"
	"backend/generation/generator/shared"
	"backend/generation/helper"
	"backend/generation/sampler"
	"backend/generation/validator"
	"backend/models"
	"math/rand/v2"
)

type TruthTableGenerator struct {
	validator validator.Validator
}

func NewTruthTableGenerator(v validator.Validator) TruthTableGenerator {
	return TruthTableGenerator{validator: v}
}

func (g TruthTableGenerator) Generate(rng *rand.Rand, prof sampler.Profile, plan sampler.Plan) (core.CandidatePools, map[string]any, error) {
	if rng == nil {
		return core.CandidatePools{}, nil, shared.ErrRngRequired
	}

	vars := shared.CanonicalVars(prof.Vars)
	attempts := 0
	for {
		if attempts >= 64 {
			return core.CandidatePools{}, nil, shared.ErrGenerationBudgetExceeded
		}
		attempts++

		// 随机生成公式
		formula := shared.RandomFormula(rng, prof)
		usedVars := shared.FilterVars(vars, formula)
		// 构建TT池
		ttPools := g.buildTruthTable(formula, usedVars)
		// 保证非平凡
		if len(ttPools.TrueSet) == 0 || len(ttPools.FalseSet) == 0 {
			continue // tautology or contradiction; regenerate
		}
		// 保证满足计划要求
		if !g.isPlanFeasible(plan, ttPools) {
			continue
		}

		pools := core.CandidatePools{TruthTable: &ttPools}
		hints := map[string]any{
			"true_count":  len(ttPools.TrueSet),
			"false_count": len(ttPools.FalseSet),
		}
		return pools, hints, nil
	}
}

// 构建TT池
func (g TruthTableGenerator) buildTruthTable(formula *core.Node, vars []string) core.TruthTablePools {
	assignments := helper.EnumerateAssignments(vars)
	trueSet := make([]string, 0, len(assignments))
	falseSet := make([]string, 0, len(assignments))

	for _, assign := range assignments {
		row := helper.AssignmentStringify(vars, assign)
		if g.validator.Eval(formula, assign) {
			trueSet = append(trueSet, row)
		} else {
			falseSet = append(falseSet, row)
		}
	}

	return core.TruthTablePools{
		Formula:  formula,
		Vars:     vars,
		TrueSet:  trueSet,
		FalseSet: falseSet,
	}
}

// 检是否满足计划要求
func (TruthTableGenerator) isPlanFeasible(plan sampler.Plan, pools core.TruthTablePools) bool {
	switch plan.QType {
	case models.QuestionTypeTrueFalse:
		// TT_EVAL_AT_ASSIGNMENT 使用，生成阶段已经排除恒真/恒假公式
		return true
	case models.QuestionTypeSingleChoice:
		switch plan.Intent {
		case "TT_TRUE_ASSIGNMENTS":
			return len(pools.TrueSet) >= 1 && len(pools.FalseSet) >= 3
		case "TT_FALSE_ASSIGNMENTS":
			return len(pools.FalseSet) >= 1 && len(pools.TrueSet) >= 3
		case "TT_EVAL_AT_ASSIGNMENT":
			return true
		default:
			return false
		}
	case models.QuestionTypeMultipleChoice:
		correct := plan.MCCorrectCount
		if correct <= 0 || correct > 4 {
			correct = 2
		}
		neededDistractor := 4 - correct
		switch plan.Intent {
		case "TT_TRUE_ASSIGNMENTS":
			return len(pools.TrueSet) >= correct && len(pools.FalseSet) >= neededDistractor
		case "TT_FALSE_ASSIGNMENTS":
			return len(pools.FalseSet) >= correct && len(pools.TrueSet) >= neededDistractor
		default:
			return false
		}
	default:
		return false
	}
}

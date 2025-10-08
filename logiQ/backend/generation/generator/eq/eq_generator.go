package eq

import (
	"backend/generation/config"
	"backend/generation/core"
	"backend/generation/generator/shared"
	"backend/generation/helper"
	"backend/generation/sampler"
	"backend/generation/validator"
	"backend/models"
	"math/rand/v2"
)

type EquivalenceGenerator struct {
	validator validator.Validator
	cfg       config.EquivalenceConfig
}

func NewEquivalenceGenerator(v validator.Validator, cfg config.EquivalenceConfig) EquivalenceGenerator {
	return EquivalenceGenerator{validator: v, cfg: cfg}
}

func (g EquivalenceGenerator) Generate(rng *rand.Rand, prof sampler.Profile, plan sampler.Plan) (core.CandidatePools, map[string]any, error) {
	if rng == nil {
		return core.CandidatePools{}, nil, shared.ErrRngRequired
	}
	attempts := 0
	vars := shared.CanonicalVars(prof.Vars)
	for {
		if attempts >= shared.MAX_ATTEMPTS {
			return core.CandidatePools{}, nil, shared.ErrGenerationBudgetExceeded
		}
		attempts++

		// 1. 随机生成目标公式 F，并确定实际使用到的变量顺序
		targetFormula := shared.RandomFormula(rng, prof)
		usedVars := shared.FilterVars(vars, targetFormula)

		//2. 利用等价/非等价规则生成候选集合
		equivCandidates := GenerateEquivalentVariants(targetFormula, rng, g.cfg, prof.EqProfile.ChainSteps)
		nonEquivCandidates := GenerateNonEquivalentVariants(targetFormula, rng, g.cfg, prof.EqProfile.ChainSteps)

		//3. run validator for each candidate, filtering by equivalence / non-equivalence
		equivPool, nonEquivPool := g.filterCandidates(targetFormula, equivCandidates, nonEquivCandidates, usedVars)

		// 4. 若候选数量不足，继续重试
		if len(equivPool) == 0 || len(nonEquivPool) == 0 {
			continue
		}

		eqPools := core.EquivalencePools{
			Target:       targetFormula,
			Vars:         usedVars,
			EquivPool:    equivPool,
			NonEquivPool: nonEquivPool,
		}

		// 5. 根据 plan.Intent / plan.QType 确认是否满足正确/干扰项数量需求
		if g.isPlanFeasible(plan, eqPools) == false {
			continue
		}

		// 6. 构造 CandidatePools 并返回
		pools := core.CandidatePools{
			Equivalence: &eqPools,
		}
		hints := map[string]any{"equiv_candidates": len(equivPool), "non_equiv_candidates": len(nonEquivPool)}
		return pools, hints, nil
	}
}

// filterCandidates
// 过滤等价候选, 非等价候选
// 删除长度过长的公式
// 返回字符串形式的公式
func (g EquivalenceGenerator) filterCandidates(target *core.Node, equivCandidates []*core.Node, nonEquivCandidates []*core.Node, vars []string) ([]string, []string) {
	equi := make([]string, 0, len(equivCandidates))
	nonEqui := make([]string, 0, len(nonEquivCandidates))

	targetStr := helper.Stringify(target)
	for _, candidate := range equivCandidates {
		candidateStr := helper.Stringify(candidate)
		if g.validator.Equivalent(target, candidate, vars) && candidateStr != targetStr && len(candidateStr) <= shared.MAX_EXPR_LENGTH {
			equi = append(equi, candidateStr)
		}
	}

	for _, candidate := range nonEquivCandidates {
		candidateStr := helper.Stringify(candidate)
		if !g.validator.Equivalent(target, candidate, vars) && len(candidateStr) <= shared.MAX_EXPR_LENGTH {
			nonEqui = append(nonEqui, candidateStr)
		}
	}

	return equi, nonEqui
}

// 检是否满足计划要求
func (EquivalenceGenerator) isPlanFeasible(plan sampler.Plan, pools core.EquivalencePools) bool {
	switch plan.QType {
	case models.QuestionTypeTrueFalse:
		// EQ_PAIR_TF 使用，生成阶段已经排除Pools为空的情况
		return true
	case models.QuestionTypeSingleChoice:
		switch plan.Intent {
		case "EQ_EQUIVALENT":
			return len(pools.EquivPool) >= 1 && len(pools.NonEquivPool) >= 3
		case "EQ_NONEQUIVALENT":
			return len(pools.NonEquivPool) >= 1 && len(pools.EquivPool) >= 3
		case "EQ_PAIR_TF":
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
		case "EQ_EQUIVALENT":
			return len(pools.EquivPool) >= correct && len(pools.NonEquivPool) >= neededDistractor
		case "EQ_NONEQUIVALENT":
			return len(pools.NonEquivPool) >= correct && len(pools.EquivPool) >= neededDistractor
		default:
			return false
		}
	default:
		return false
	}
}

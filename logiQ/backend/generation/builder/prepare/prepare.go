package prepare

import (
	"backend/generation/core"
	"backend/generation/helper"
	"backend/generation/sampler"
	"backend/models"
	"fmt"
	"math/rand/v2"
)

type PrepareResult struct {
	PromptData map[string]string
	TFAnswer   bool // only valid if QType is TrueFalse
}

// PreparePromptData 预处理环节
// 准备 Prompt 所需的模板映射数据
// 以及（如果适用）True/False 问题的答案。
func PreparePromptData(plan sampler.Plan, pools core.CandidatePools, rng *rand.Rand) (PrepareResult, error) {
	data := make(map[string]string)
	var tfAnswer bool

	switch plan.Category {
	case models.QuestionCategoryTruthTable:
		data["F"] = helper.Stringify(pools.TruthTable.Formula)
		if plan.QType == models.QuestionTypeTrueFalse {
			assignment, isTrue, err := sampleTruthTableTF(pools.TruthTable, rng)
			if err != nil {
				return PrepareResult{}, err
			}
			data["alpha"] = assignment
			tfAnswer = isTrue
		}
	case models.QuestionCategoryEquivalence:
		data["T"] = helper.Stringify(pools.Equivalence.Target)
		if plan.QType == models.QuestionTypeTrueFalse {
			partner, isEquiv, err := sampleEquivalenceTF(pools.Equivalence, rng)
			if err != nil {
				return PrepareResult{}, err
			}
			data["G"] = partner
			tfAnswer = isEquiv
		}
	case models.QuestionCategoryInference:
		data["Premises"] = pools.Inference.Premises
		if plan.QType == models.QuestionTypeTrueFalse {
			conclusion, isValid, err := sampleInferenceTF(pools.Inference, rng)
			if err != nil {
				return PrepareResult{}, err
			}
			data["Conclusion"] = conclusion
			tfAnswer = isValid
		}
	default:
		return PrepareResult{}, fmt.Errorf("service: unsupported category %s", plan.Category)
	}

	return PrepareResult{
		PromptData: data,
		TFAnswer:   tfAnswer,
	}, nil

}

func sampleTruthTableTF(pools *core.TruthTablePools, rng *rand.Rand) (string, bool, error) {
	if rng == nil {
		return "", false, fmt.Errorf("service: rng must not be nil")
	}
	if pools == nil {
		return "", false, fmt.Errorf("service: nil truth-table pools")
	}

	total := len(pools.TrueSet) + len(pools.FalseSet)
	if total == 0 {
		return "", false, fmt.Errorf("service: truth-table pools empty")
	}

	idx := rng.IntN(total)
	if idx < len(pools.TrueSet) {
		return pools.TrueSet[idx], true, nil
	}
	return pools.FalseSet[idx-len(pools.TrueSet)], false, nil
}

func sampleEquivalenceTF(pools *core.EquivalencePools, rng *rand.Rand) (string, bool, error) {
	if rng == nil {
		return "", false, fmt.Errorf("service: rng must not be nil")
	}
	if pools == nil {
		return "", false, fmt.Errorf("service: nil equivalence pools")
	}

	total := len(pools.EquivPool) + len(pools.NonEquivPool)
	if total == 0 {
		return "", false, fmt.Errorf("service: equivalence pools empty")
	}

	idx := rng.IntN(total)
	if idx < len(pools.EquivPool) {
		return pools.EquivPool[idx], true, nil
	}
	return pools.NonEquivPool[idx-len(pools.EquivPool)], false, nil
}

func sampleInferenceTF(pools *core.InferencePools, rng *rand.Rand) (string, bool, error) {
	if rng == nil {
		return "", false, fmt.Errorf("service: rng must not be nil")
	}
	if pools == nil {
		return "", false, fmt.Errorf("service: nil inference pools")
	}

	total := len(pools.ValidConclusions) + len(pools.InvalidConclusions)
	if total == 0 {
		return "", false, fmt.Errorf("service: inference pools empty")
	}

	idx := rng.IntN(total)
	if idx < len(pools.ValidConclusions) {
		return pools.ValidConclusions[idx], true, nil
	}
	return pools.InvalidConclusions[idx-len(pools.ValidConclusions)], false, nil
}

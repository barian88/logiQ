package choice

import (
	"backend/generation/config"
	"backend/generation/core"
	"backend/generation/sampler"
	"backend/models"
	"errors"
	"fmt"
	"math/rand/v2"
	"sort"
)

// Choice holds the shuffled options together with the indices of correct answers.
type Choice struct {
	Options        []string
	CorrectIndexes []int
}

var (
	ErrUnsupportedIntent = errors.New("choice: unsupported intent for given category")
	ErrMissingPool       = errors.New("choice: required candidate pool not available")
	ErrNilRNG            = errors.New("choice: rng must not be nil")
)

// Builder draws options from candidate pools according to the intent specification.
type Builder struct{}

func NewBuilder() Builder { return Builder{} }

// Params bundles inputs required to assemble choices.
type Params struct {
	Plan     sampler.Plan
	Intent   config.IntentSpec
	Pools    core.CandidatePools
	TFAnswer bool // 对于 trueFalse 题型必须提供正确答案
}

// BuildChoice selects options from CandidatePools following the plan/intent contract.
func (Builder) BuildChoice(params Params, rng *rand.Rand) (Choice, error) {
	if rng == nil {
		return Choice{}, ErrNilRNG
	}

	switch params.Plan.QType {
	case models.QuestionTypeSingleChoice:
		return buildSC(params.Plan, params.Intent, params.Pools, rng)
	case models.QuestionTypeMultipleChoice:
		return buildMC(params.Plan, params.Intent, params.Pools, rng)
	case models.QuestionTypeTrueFalse:
		return buildTF(params)
	default:
		return Choice{}, fmt.Errorf("choice: unsupported question type %s", params.Plan.QType)
	}
}

func buildSC(plan sampler.Plan, intent config.IntentSpec, pools core.CandidatePools, rng *rand.Rand) (Choice, error) {
	// 单选：从配置映射到的正确池抽 1、从干扰池抽 3
	correctPool, distractorPool, err := resolvePools(plan, intent.PoolMapping.SC, pools)
	if err != nil {
		return Choice{}, err
	}

	correct := sampleUnique(correctPool, 1, rng)
	distractors := sampleUnique(distractorPool, 3, rng)

	options, indexes := combineAndShuffle(correct, distractors, rng)
	return Choice{Options: options, CorrectIndexes: indexes}, nil
}

func buildMC(plan sampler.Plan, intent config.IntentSpec, pools core.CandidatePools, rng *rand.Rand) (Choice, error) {
	// 多选：从正确池抽 plan.MCCorrectCount 个，其余补足 4 个选项
	correctPool, distractorPool, err := resolvePools(plan, intent.PoolMapping.MC, pools)
	if err != nil {
		return Choice{}, err
	}
	// 确保正确选项数在合理范围内
	k := plan.MCCorrectCount
	if k <= 0 {
		k = 2
	}
	if k > 4 {
		k = 4
	}
	correct := sampleUnique(correctPool, k, rng)
	distractors := sampleUnique(distractorPool, 4-k, rng)

	options, indexes := combineAndShuffle(correct, distractors, rng)
	return Choice{Options: options, CorrectIndexes: indexes}, nil
}

func buildTF(params Params) (Choice, error) {
	// 统一的 True/False 选项，由调用方提供正确答案
	options := []string{"True", "False"}
	correctIndex := 1
	if params.TFAnswer {
		correctIndex = 0
	}
	return Choice{Options: options, CorrectIndexes: []int{correctIndex}}, nil
}

func resolvePools(plan sampler.Plan, mapping *config.PoolMapping, pools core.CandidatePools) ([]string, []string, error) {
	if mapping == nil {
		return nil, nil, ErrUnsupportedIntent
	}

	switch plan.Category {
	case models.QuestionCategoryTruthTable:
		if pools.TruthTable == nil {
			return nil, nil, ErrMissingPool
		}
		return selectPool(mapping, map[string][]string{
			"TrueSet":  pools.TruthTable.TrueSet,
			"FalseSet": pools.TruthTable.FalseSet,
		})
	case models.QuestionCategoryEquivalence:
		if pools.Equivalence == nil {
			return nil, nil, ErrMissingPool
		}
		return selectPool(mapping, map[string][]string{
			"EquivPool":    pools.Equivalence.EquivPool,
			"NonEquivPool": pools.Equivalence.NonEquivPool,
		})
	case models.QuestionCategoryInference:
		if pools.Inference == nil {
			return nil, nil, ErrMissingPool
		}
		return selectPool(mapping, map[string][]string{
			"ValidConclusions":   pools.Inference.ValidConclusions,
			"InvalidConclusions": pools.Inference.InvalidConclusions,
		})
	default:
		return nil, nil, ErrUnsupportedIntent
	}
}

func selectPool(mapping *config.PoolMapping, pools map[string][]string) ([]string, []string, error) {
	// 依据 Intent 配置中的字符串名称映射到真实候选切片
	correct, ok := pools[mapping.Correct]
	if !ok {
		return nil, nil, fmt.Errorf("choice: unknown pool %s", mapping.Correct)
	}
	distractor, ok := pools[mapping.Distractor]
	if !ok {
		return nil, nil, fmt.Errorf("choice: unknown pool %s", mapping.Distractor)
	}
	return correct, distractor, nil
}

func sampleUnique(pool []string, n int, rng *rand.Rand) []string {
	// 随机抽取 n 个互不重复元素；不足时取全部
	if n <= 0 {
		return nil
	}
	if len(pool) == 0 {
		return nil
	}
	if n > len(pool) {
		n = len(pool)
	}

	perm := rng.Perm(len(pool))
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = pool[perm[i]]
	}
	return out
}

func combineAndShuffle(correct, distractors []string, rng *rand.Rand) ([]string, []int) {
	// 合并正确/干扰项并打乱，返回乱序选项与正确索引
	combined := append(append([]string{}, correct...), distractors...)
	perm := rng.Perm(len(combined))

	shuffled := make([]string, len(combined))
	correctIndexes := make([]int, 0, len(correct))
	targets := make(map[string]int, len(correct))
	for _, c := range correct {
		targets[c]++
	}

	for i, idx := range perm {
		val := combined[idx]
		shuffled[i] = val
		if count, ok := targets[val]; ok {
			correctIndexes = append(correctIndexes, i)
			if count == 1 {
				delete(targets, val)
			} else {
				targets[val] = count - 1
			}
		}
	}

	sort.Ints(correctIndexes)
	return shuffled, correctIndexes
}

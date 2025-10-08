package sampler

import (
	"backend/generation/config"
	"backend/generation/core"
	"backend/generation/helper"
	"backend/models"
	"fmt"
	"math/rand/v2"
	"strings"
)

// Plan 计划（Plan）定义了题目的高层次要求，会被多个层使用。
type Plan struct {
	Difficulty     models.QuestionDifficulty
	Category       models.QuestionCategory
	QType          models.QuestionType
	Intent         string
	MCCorrectCount int
}

// EqProfile 等价题型的配置文件
type EqProfile struct {
	ChainSteps int
}

// InfProfile 推理题型的配置文件
type InfProfile struct {
	ChainSteps int
}

// Profile 配置文件（Profile）定义了题目生成的参数，Generator 负责使用。
type Profile struct {
	Vars       int
	MaxDepth   int
	AllowedOps []core.NodeKind
	EqProfile  EqProfile
	InfProfile InfProfile
}

// SampleResult 采样结果，包含 Plan 和 Profile。
type SampleResult struct {
	Plan    Plan
	Profile Profile
}

type Sampler struct {
	cfg config.AppConfig
}

func NewSampler(cfg config.AppConfig) Sampler {
	return Sampler{cfg: cfg}
}

// Sample 根据(三个可选参数)题目难度、题型和类别，生成一个采样结果（SampleResult）。
func (s Sampler) Sample(rng *rand.Rand, difficulty models.QuestionDifficulty, qType models.QuestionType, category models.QuestionCategory) (SampleResult, error) {
	// 第一步抽样 Plan
	plan, err := s.samplePlan(rng, difficulty, qType, category)
	if err != nil {
		return SampleResult{}, err
	}

	// 第二步抽样 Profile
	profile, err := s.sampleProfile(rng, plan)
	if err != nil {
		return SampleResult{}, err
	}
	return SampleResult{Plan: plan, Profile: profile}, nil
}

// samplePlan 抽样 Plan
func (s Sampler) samplePlan(rng *rand.Rand, difficulty models.QuestionDifficulty, qType models.QuestionType, category models.QuestionCategory) (Plan, error) {
	cfg := s.cfg

	// 难度
	planDifficulty := difficulty
	// 如果没有指定难度，则根据权重随机选择一个难度
	if planDifficulty == "" {
		planDifficulty = helper.SampleWeighted(cfg.Planner.DifficultyWeights, rng)
	}
	// 类别
	planCategory := category
	// 如果没有指定类别，则根据难度对应的类别权重随机选择一个类别
	if planCategory == "" {
		catDist := cfg.Planner.CategoryWeights[planDifficulty]
		planCategory = helper.SampleWeighted(catDist, rng)
	}
	// 题型
	planQType := qType
	if planQType == "" {
		typeDist := cfg.Planner.TypeWeights[planDifficulty]
		planQType = helper.SampleWeighted(typeDist, rng)
	}

	// Intent
	intentDistByDifficulty, ok := cfg.Planner.IntentWeights[planCategory]
	if !ok {
		return Plan{}, fmt.Errorf("sampler: no intent weights for category %s", planCategory)
	}
	intentDist := intentDistByDifficulty[planDifficulty]
	// 过滤掉权重为0的意图，以及不支持当前题型的意图
	// 比如说type抽到了True/False，但是抽到了true assignment的intent，它只支持SC和MC，就不能用
	filtered := make(map[string]float64)
	for name, weight := range intentDist {
		spec, ok := cfg.Intents[name]
		if !ok {
			continue
		}
		if len(spec.Templates.TemplatesFor(planQType)) == 0 {
			continue
		}
		filtered[name] = weight
	}
	if len(filtered) == 0 {
		return Plan{}, fmt.Errorf("sampler: no intents available for category=%s difficulty=%s type=%s", planCategory, planDifficulty, planQType)
	}
	planIntent := helper.SampleWeighted(filtered, rng)

	plan := Plan{
		Difficulty: planDifficulty,
		Category:   planCategory,
		QType:      planQType,
		Intent:     planIntent,
	}

	// 多选题的正确选项数
	if planQType == models.QuestionTypeMultipleChoice {
		countDist := cfg.Planner.MCCorrectCountDist[planDifficulty]
		plan.MCCorrectCount = helper.SampleWeighted(countDist, rng)
	}
	return plan, nil
}

// sampleProfile 根据 Plan 抽样 Profile
func (s Sampler) sampleProfile(rng *rand.Rand, plan Plan) (Profile, error) {
	cfg := s.cfg

	vars := helper.SampleWeighted(cfg.DifficultyProfiles[plan.Difficulty].VarsDist, rng)
	maxDepth := helper.SampleWeighted(cfg.DifficultyProfiles[plan.Difficulty].DepthDist, rng)
	allowedOpsOrig := cfg.DifficultyProfiles[plan.Difficulty].AllowedOps
	allowedOps := make([]core.NodeKind, 0, len(allowedOpsOrig))
	for _, opStr := range allowedOpsOrig {
		op, ok := config.OpNameToKind[strings.ToUpper(opStr)]
		if !ok {
			return Profile{}, fmt.Errorf("sampler: unsupported operator %s", opStr)
		}
		allowedOps = append(allowedOps, op)
	}
	profile := Profile{
		Vars:       vars,
		MaxDepth:   maxDepth,
		AllowedOps: allowedOps,
	}
	// Equivalence 题型，额外采样链长
	if plan.Category == models.QuestionCategoryEquivalence {
		diffCfg, ok := cfg.Equivalence.Difficulty[plan.Difficulty]
		if !ok {
			return Profile{}, fmt.Errorf("sampler: equivalence difficulty %s not configured", plan.Difficulty)
		}
		chainSteps := helper.SampleWeighted(diffCfg.ChainStepsDist, rng)
		eqProfile := EqProfile{
			ChainSteps: chainSteps,
		}
		profile.EqProfile = eqProfile
	}
	// Inference 题型，额外采样链长和冗余前提数量
	if plan.Category == models.QuestionCategoryInference {
		diffCfg, ok := cfg.Inference.Difficulty[plan.Difficulty]
		if !ok {
			return Profile{}, fmt.Errorf("sampler: inference difficulty %s not configured", plan.Difficulty)
		}
		chainSteps := helper.SampleWeighted(diffCfg.ChainStepsDist, rng)
		infProfile := InfProfile{
			ChainSteps: chainSteps,
		}
		profile.InfProfile = infProfile
	}
	return profile, nil
}

package inf

import (
	"backend/generation/config"
	"backend/generation/core"
	"backend/generation/generator/shared"
	"backend/generation/sampler"
	"backend/generation/validator"
	"backend/models"
	"fmt"
	"math/rand/v2"
)

type InferenceGenerator struct {
	cfg       config.InferenceConfig
	validator validator.Validator
}

func NewInferenceGenerator(validator validator.Validator, cfg config.InferenceConfig) InferenceGenerator {
	return InferenceGenerator{
		cfg:       cfg,
		validator: validator,
	}
}

// Generate 生成一个推理题目
func (g InferenceGenerator) Generate(rng *rand.Rand, prof sampler.Profile, plan sampler.Plan) (core.CandidatePools, map[string]any, error) {
	if rng == nil {
		return core.CandidatePools{}, nil, shared.ErrRngRequired
	}
	attempts := 0
	for {
		if attempts >= 64 {
			return core.CandidatePools{}, nil, shared.ErrGenerationBudgetExceeded
		}
		attempts++

		pair, err := g.selectTemplatePairRandomly(g.cfg.TemplatePairs, prof.InfProfile.ChainSteps, rng)
		if err != nil {
			return core.CandidatePools{}, nil, err
		}
		// 解析前提和结论模板
		premises, validCons, inValidCons, err := extractPremisesAndConclusions(*pair)
		if err != nil {
			return core.CandidatePools{}, nil, err
		}
		// 获取占位符的绑定
		bindings, err := g.prepareBindings(rng, pair.Slots)
		if err != nil {
			return core.CandidatePools{}, nil, err
		}
		// 替换占位符
		premisesInst, validConInst, inValidConInst := instantiateTemplatePair(premises, validCons, inValidCons, bindings)

		// 跳过验证的步骤，因为模板已经保证了正确性，验证inf性能开销较大

		// 对前提和结论进行等价变换，拓展，去重等  增加多样性
		finalPremise, transValid, transInvalid := ExpandAndTransform(premisesInst, validConInst, inValidConInst, rng)

		// 验证结论的正确性
		usedVars := collectVars(finalPremise)
		finalValid, finalInvalid := g.validateInference(finalPremise, transValid, transInvalid, usedVars)

		// 验证生成的题目是否满足计划要求
		if g.isPlanFeasible(plan, finalValid, finalInvalid) == false {
			continue
		}

		// stringify
		premiseStr, validStrs, inValidStrs := stringifyInstances(finalPremise, finalValid, finalInvalid)

		// 构建pools
		infPools := core.InferencePools{
			Premises:           premiseStr,
			ValidConclusions:   validStrs,
			InvalidConclusions: inValidStrs,
			Vars:               usedVars, // 只考虑前提中的变量，结论不会引入新变量
		}

		return core.CandidatePools{
			Inference: &infPools,
		}, nil, nil
	}

}

// selectTemplatePairRandomly 从模板对列表中，随机选择一个满足链条长度要求的模板对
func (g InferenceGenerator) selectTemplatePairRandomly(pairs []config.InferenceTemplatePair, chainSteps int, rng *rand.Rand) (*config.InferenceTemplatePair, error) {
	candidates := make([]config.InferenceTemplatePair, 0)
	for _, pair := range pairs {
		if pair.ChainSteps == chainSteps {
			candidates = append(candidates, pair)
		}
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no template pair found for chain steps: %d", chainSteps)
	}
	selected := candidates[rng.IntN(len(candidates))]
	return &selected, nil
}

// isPlanFeasible 检查生成的前提和结论是否满足计划要求
func (g InferenceGenerator) isPlanFeasible(plan sampler.Plan, validCons []*core.Node, inValidCons []*core.Node) bool {
	switch plan.QType {
	case models.QuestionTypeTrueFalse:
		// INF_VALIDITY_TF 使用，模板配置已经排除Pools为空的情况
		return true
	case models.QuestionTypeSingleChoice:
		switch plan.Intent {
		case "INF_DERIVABLE":
			return len(validCons) >= 1 && len(inValidCons) >= 3
		case "INF_UNDERIVABLE":
			return len(inValidCons) >= 1 && len(validCons) >= 3
		case "INF_VALIDITY_TF":
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
		case "INF_DERIVABLE":
			return len(validCons) >= correct && len(inValidCons) >= neededDistractor
		case "INF_UNDERIVABLE":
			return len(inValidCons) >= correct && len(validCons) >= neededDistractor
		default:
			return false
		}
	default:
		return false
	}

}

func (g InferenceGenerator) validateInference(premises, validCon, invalidCon []*core.Node, vars []string) ([]*core.Node, []*core.Node) {
	validatedValid := make([]*core.Node, 0, len(validCon))
	validatedInvalid := make([]*core.Node, 0, len(invalidCon))
	for _, node := range validCon {
		if g.validator.Derivable(premises, node, vars) {
			validatedValid = append(validatedValid, node)
		}
	}
	for _, node := range invalidCon {
		if !g.validator.Derivable(premises, node, vars) {
			validatedInvalid = append(validatedInvalid, node)
		}
	}
	return validatedValid, validatedInvalid
}

package service

import (
	"backend/generation/assembler"
	"backend/generation/builder/choice"
	"backend/generation/builder/prepare"
	"backend/generation/builder/prompt"
	"backend/generation/config"
	"backend/generation/generator"
	"backend/generation/generator/eq"
	"backend/generation/generator/inf"
	"backend/generation/generator/tt"
	"backend/generation/sampler"
	"backend/generation/validator"
	"backend/models"
	"fmt"
	"math/rand/v2"
	"time"
)

type Service struct {
	cfg           config.AppConfig
	sampler       sampler.Sampler
	generators    map[models.QuestionCategory]generator.Generator
	promptBuilder prompt.Builder
	choiceBuilder choice.Builder
	assembler     assembler.Assembler
}

func NewService() Service {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	//
	generators := make(map[models.QuestionCategory]generator.Generator)
	defaultValidator := validator.NewDefaultValidator()
	generators[models.QuestionCategoryTruthTable] = tt.NewTruthTableGenerator(defaultValidator)
	generators[models.QuestionCategoryEquivalence] = eq.NewEquivalenceGenerator(defaultValidator, cfg.Equivalence)
	generators[models.QuestionCategoryInference] = inf.NewInferenceGenerator(defaultValidator, cfg.Inference)
	return Service{
		cfg:           cfg,
		sampler:       sampler.NewSampler(cfg),
		generators:    generators,
		promptBuilder: prompt.NewBuilder(),
		choiceBuilder: choice.NewBuilder(),
		assembler:     assembler.NewAssembler(),
	}
}

func (s Service) GenerateQuestion(num int, category models.QuestionCategory, difficulty models.QuestionDifficulty, qType models.QuestionType) ([]models.Question, error) {

	if num <= 0 {
		return nil, nil
	}

	rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	questionList := make([]models.Question, 0, num)
	const maxAttemptsPerQuestion = 5

	// 循环生成题目，直到达到所需数量
	for len(questionList) < num {
		var lastErr error
		succeeded := false
		// 每次生成题目，最多尝试maxAttemptsPerQuestion次，达到则视为用户的需求无法满足，放弃生成并返回，避免卡死
		for attempt := 0; attempt < maxAttemptsPerQuestion; attempt++ {
			// 1. sample
			sampleResult, err := s.sampler.Sample(rng, difficulty, qType, category)
			if err != nil {
				lastErr = err
				continue
			}
			plan := sampleResult.Plan
			profile := sampleResult.Profile

			// 2. generate
			generatorSpec, ok := s.generators[plan.Category]
			if !ok {
				lastErr = fmt.Errorf("service: no generator registered for category %s", plan.Category)
				continue
			}
			candidatePools, _, err := generatorSpec.Generate(rng, profile, plan)
			if err != nil {
				lastErr = err
				continue
			}

			// 3. build prompt
			// 获取intentSpec
			intentSpec, ok := s.cfg.Intents[plan.Intent]
			if !ok {
				lastErr = fmt.Errorf("service: intent %s not found", plan.Intent)
				continue
			}
			buildCtx, err := prepare.PreparePromptData(plan, candidatePools, rng)
			if err != nil {
				lastErr = err
				continue
			}
			var promptParam = prompt.Params{
				Intent: intentSpec,
				QType:  plan.QType,
				Data:   buildCtx.PromptData,
			}
			promptRes, err := s.promptBuilder.BuildPrompt(promptParam, rng)
			if err != nil {
				lastErr = err
				continue
			}

			// 4. build choice
			var choiceParam = choice.Params{
				Plan:     plan,
				Intent:   intentSpec,
				Pools:    candidatePools,
				TFAnswer: buildCtx.TFAnswer,
			}

			choiceRes, err := s.choiceBuilder.BuildChoice(choiceParam, rng)
			if err != nil {
				lastErr = err
				continue
			}

			// 5. assemble
			question := s.assembler.Assemble(plan, promptRes, choiceRes)
			questionList = append(questionList, question)
			succeeded = true
			break
		}

		if !succeeded {
			return nil, fmt.Errorf("service: failed to generate question after %d attempts: %w", maxAttemptsPerQuestion, lastErr)
		}
	}

	// 6. 返回
	return questionList, nil
}

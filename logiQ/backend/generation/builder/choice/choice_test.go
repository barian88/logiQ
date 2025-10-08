package choice

import (
	"backend/generation/config"
	"backend/generation/core"
	"backend/generation/generator/tt"
	"backend/generation/helper"
	"backend/generation/sampler"
	"backend/generation/validator"
	"backend/models"
	"math/rand/v2"
	"strings"
	"testing"
)

var (
	rng  *rand.Rand      = rand.New(rand.NewPCG(1, 5))
	prof sampler.Profile = sampler.Profile{
		Vars:       3,
		MaxDepth:   2,
		AllowedOps: []core.NodeKind{core.Not, core.And, core.Or, core.Impl, core.Iff},
	}
	ttGen tt.TruthTableGenerator = tt.NewTruthTableGenerator(validator.NewDefaultValidator())
	plan                         = sampler.Plan{Category: models.QuestionCategoryTruthTable,

		QType: models.QuestionTypeMultipleChoice, Intent: "TT_FALSE_ASSIGNMENTS", MCCorrectCount: 2}
)

func TestBuildChoice(t *testing.T) {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	library := (cfg.Intents)

	// 生成一个题目
	pools, _, err := ttGen.Generate(rng, prof, plan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("question formula: %s", helper.Stringify(pools.TruthTable.Formula))
	intentSpec, ok := library[plan.Intent]
	if !ok {
		t.Fatalf("intent not found")
	}
	var param = Params{
		Plan:   plan,
		Intent: intentSpec,
		Pools:  pools,
		//TFAnswer: false,
	}
	choiceBuilder := NewBuilder()
	choice, err := choiceBuilder.BuildChoice(param, rng)
	if err != nil {
		t.Fatalf("failed to build choice: %v", err)
	}
	t.Logf("Options: %s", strings.Join(choice.Options, "\n"))
	t.Logf("CorrectIndexes: %v", choice.CorrectIndexes)

}

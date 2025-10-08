package inf

import (
	"backend/generation/config"
	"backend/generation/core"
	"backend/generation/sampler"
	"backend/generation/validator"
	"backend/models"
	"math/rand/v2"
	"testing"
)

var (
	rng         *rand.Rand = rand.New(rand.NewPCG(1, 4))
	myValidator            = validator.NewDefaultValidator()
	prof                   = sampler.Profile{
		Vars:       3,
		MaxDepth:   3,
		AllowedOps: []core.NodeKind{core.Not, core.And, core.Or, core.Impl, core.Iff},
		InfProfile: sampler.InfProfile{
			ChainSteps: 2,
		},
	}
	plan = sampler.Plan{QType: models.QuestionTypeSingleChoice, Intent: "INF_DERIVABLE", MCCorrectCount: 0}
)

func TestInfGenerator(t *testing.T) {

	cfg, _ := config.LoadConfig()
	infG := NewInferenceGenerator(myValidator, cfg.Inference)

	pools, _, err := infG.Generate(rng, prof, plan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	premiseStr := pools.Inference.Premises
	t.Logf("premises: %v", premiseStr)
	for _, v := range pools.Inference.ValidConclusions {
		t.Logf("Valid conclusion: %s \n", v)
	}
	for _, iv := range pools.Inference.InvalidConclusions {
		t.Logf("Invalid conclusion: %s \n", iv)
	}

}

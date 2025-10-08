package eq

import (
	"backend/generation/config"
	"backend/generation/core"
	"backend/generation/helper"
	"backend/generation/sampler"
	"backend/generation/validator"
	"backend/models"
	"math/rand/v2"
	"testing"
)

var (
	rng  *rand.Rand      = rand.New(rand.NewPCG(1, 8))
	prof sampler.Profile = sampler.Profile{
		Vars:       4,
		MaxDepth:   4,
		AllowedOps: []core.NodeKind{core.Not, core.And, core.Or, core.Impl, core.Iff},
		EqProfile: sampler.EqProfile{
			ChainSteps: 2,
		},
	}
	plan = sampler.Plan{QType: models.QuestionTypeMultipleChoice, Intent: "EQ_EQUIVALENT", MCCorrectCount: 2}
)

func TestGenerate(t *testing.T) {
	validator := validator.NewDefaultValidator()
	appCfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	generator := NewEquivalenceGenerator(validator, appCfg.Equivalence)
	pools, _, err := generator.Generate(rng, prof, plan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pools.Equivalence == nil {
		t.Fatal("expected equivalence pools, got nil")
	}
	eqPools := pools.Equivalence
	t.Logf("Target formula: %s", helper.Stringify(eqPools.Target))
	t.Logf("Used vars: %v", eqPools.Vars)
	// 遍历等价候选并输出
	for i, candidate := range eqPools.EquivPool {
		t.Logf("Equiv Candidate %d: %s", i+1, candidate)
		// 输出长度
		t.Logf("Length: %d", len(candidate))
	}
	// 遍历非等价候选并输出
	for i, candidate := range eqPools.NonEquivPool {
		t.Logf("Non-Equiv Candidate %d: %s", i+1, candidate)
	}

}

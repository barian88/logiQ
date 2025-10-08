// backend/generation/generator/generator_tt_test.go
package tt

import (
	"backend/generation/core"
	"backend/generation/generator/shared"
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
		Vars:       2,
		MaxDepth:   2,
		AllowedOps: []core.NodeKind{core.Not, core.And, core.Or, core.Impl, core.Iff},
	}
	ttGen TruthTableGenerator = NewTruthTableGenerator(validator.NewDefaultValidator())
	plan                      = sampler.Plan{QType: models.QuestionTypeMultipleChoice, Intent: "TT_TRUE_ASSIGNMENTS", MCCorrectCount: 2}
)

func TestRandomFormula(t *testing.T) {
	root := shared.RandomFormula(rng, prof)
	if root == nil {
		t.Fatal("expected a formula, got nil")
	}

	//查看生成的公式
	t.Logf("\n Generated formula: %s", helper.Stringify(root))
}

func TestTTPool(t *testing.T) {
	root := shared.RandomFormula(rng, prof)
	ttTable := ttGen.buildTruthTable(root, shared.CanonicalVars(prof.Vars))
	//查看生成的真值表
	// t.Logf("%s", strings.Join(ttTable.Vars, ","))
	t.Logf("%s", strings.Join(ttTable.TrueSet, "\n"))
	t.Logf("%s", strings.Join(ttTable.FalseSet, "\n"))
}

func TestIsPlanFeasible(t *testing.T) {
	for attempts := 0; attempts < 64; attempts++ {
		root := shared.RandomFormula(rng, prof)
		ttTable := ttGen.buildTruthTable(root, shared.CanonicalVars(prof.Vars))
		if ttGen.isPlanFeasible(plan, ttTable) {
			return
		}
	}
	t.Fatal("failed to find feasible formula after 64 attempts")
}

func TestGenerate(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	pools, hints, err := ttGen.Generate(rng, prof, plan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pools.TruthTable == nil {
		t.Fatal("expected truth table pools, got nil")
	}
	ttTable := pools.TruthTable
	t.Logf("Generated formula: %s", helper.Stringify(ttTable.Formula))
	t.Logf("Generated Truth Table Pools: \n Vars: %s \n TrueSet: %s \n FalseSet: %s", strings.Join(ttTable.Vars, ","), strings.Join(ttTable.TrueSet, "\n"), strings.Join(ttTable.FalseSet, "\n"))
	t.Logf("Hints: %v", hints)
}

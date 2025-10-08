// // 先实现 TruthTableGenerator（MVP）：
// 	•	输入：prof（Vars/MaxDepth/AllowedOps）、plan.MCCorrectCount
// 	•	产出：TruthTablePools{Formula, Vars, TrueSet, FalseSet}
// 	•	保证：
// 	•	TrueSet 与 FalseSet 非空；
// 	•	能满足 SC（1 真 + 3 假）或 MC（k 真 + (4-k) 假）；否则返回 ErrInsufficientCandidates（由 orchestrator 重试/回退）

// Equivalence/Inference 的 Generator 先留接口，下一步再具体化。

package generator

import (
	"backend/generation/core"
	"backend/generation/sampler"
	"math/rand/v2"
)

type Generator interface {
	Generate(rng *rand.Rand, prof sampler.Profile, plan sampler.Plan) (core.CandidatePools, map[string]any /*hints*/, error)
}

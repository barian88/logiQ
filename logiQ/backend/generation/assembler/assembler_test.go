package assembler

import (
	"backend/generation/builder/choice"
	"backend/generation/builder/prompt"
	"backend/generation/sampler"
	"backend/models"
	"encoding/json"
	"testing"
)

func TestAssembler_Assemble(t *testing.T) {
	var (
		prompt = prompt.Prompt{
			Text: "which is true?",
		}
		chioice = choice.Choice{
			Options:        []string{"p and q", "p or q", "not p"},
			CorrectIndexes: []int{0},
		}
		plan = sampler.Plan{
			Difficulty:     models.QuestionDifficultyEasy,
			Category:       models.QuestionCategoryTruthTable,
			QType:          models.QuestionTypeSingleChoice,
			Intent:         "TT_TRUE_ASSIGNMENTS",
			MCCorrectCount: 1,
		}
	)
	assembler := NewAssembler()
	question := assembler.Assemble(plan, prompt, chioice)
	// // 把question转换成json打印出来
	jsonBytes, err := json.Marshal(question)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	t.Logf("Generated Question JSON: %s", string(jsonBytes))
}

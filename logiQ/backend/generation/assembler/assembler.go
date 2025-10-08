package assembler

import (
	"backend/generation/builder/choice"
	"backend/generation/builder/prompt"
	"backend/generation/sampler"
	"backend/models"
)

type Assembler struct{}

func NewAssembler() Assembler {
	return Assembler{}
}

// Assemble 把 prompt 和 choice 组装成最终的 question
func (a Assembler) Assemble(plan sampler.Plan, prompt prompt.Prompt, choice choice.Choice) models.Question {

	return models.Question{
		Difficulty:         plan.Difficulty,
		Category:           plan.Category,
		Type:               plan.QType,
		QuestionText:       prompt.Text,
		Options:            choice.Options,
		CorrectAnswerIndex: choice.CorrectIndexes,
		IsActive:           true,
	}
}

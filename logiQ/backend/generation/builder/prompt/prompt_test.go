package prompt

import (
	"backend/generation/config"
	"backend/models"
	"math/rand/v2"
	"testing"
)

func TestPromptBuilder_BuildQuestionText(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	library := cfg.Intents

	intentName := "TT_TRUE_ASSIGNMENTS"
	intentSpec, ok := library[intentName]
	if !ok {
		t.Fatalf("intent not found")
	}
	promptBuilder := NewBuilder()

	// 组装一个data
	var data = make(map[string]string)
	// 假如生成的F是 "p and q"
	data["F"] = "p and q"

	rng := rand.New(rand.NewPCG(1, 2))

	prompt, err := promptBuilder.BuildPrompt(Params{
		Intent: intentSpec,
		QType:  models.QuestionTypeSingleChoice,
		Data:   data,
	}, rng)
	if err != nil {
		return
	}
	t.Logf("Generated Prompt: %s", prompt.Text)

}

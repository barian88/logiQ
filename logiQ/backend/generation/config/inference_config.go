package config

import "backend/models"

type ExprPatternSpec struct {
	Kind  string           `yaml:"kind"`
	Name  string           `yaml:"name,omitempty"`
	Left  *ExprPatternSpec `yaml:"left,omitempty"`
	Right *ExprPatternSpec `yaml:"right,omitempty"`
}

type InferenceTemplateSlot struct {
	Type string `yaml:"type"`
}

type InferenceTemplateConclusion struct {
	Name        string          `yaml:"name"`
	Description string          `yaml:"description,omitempty"`
	Conclusion  ExprPatternSpec `yaml:"conclusion"`
}

type InferenceTemplatePair struct {
	Name        string                           `yaml:"name"`
	Description string                           `yaml:"description"`
	ChainSteps  int                              `yaml:"chain_steps"`
	Slots       map[string]InferenceTemplateSlot `yaml:"slots"`
	Premises    []ExprPatternSpec                `yaml:"premises"`
	Valid       []InferenceTemplateConclusion    `yaml:"valid"`
	Invalid     []InferenceTemplateConclusion    `yaml:"invalid,omitempty"`
}

type InferenceDifficultyConfig struct {
	ChainStepsDist map[int]float64 `yaml:"chain_steps_dist"`
}

type SlotFillerSpec map[string]float64

type InferenceConfig struct {
	TemplatePairs []InferenceTemplatePair                                 `yaml:"template_pairs"`
	SlotFillers   map[string]SlotFillerSpec                               `yaml:"slot_fillers"`
	Difficulty    map[models.QuestionDifficulty]InferenceDifficultyConfig `yaml:"difficulty"`
}

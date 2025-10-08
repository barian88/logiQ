package config

import "backend/models"

type DifficultyProfile struct {
	VarsDist   map[int]float64 `yaml:"vars_dist"`
	DepthDist  map[int]float64 `yaml:"depth_dist"`
	AllowedOps []string        `yaml:"allowed_ops"`
}

type PlannerWeights struct {
	DifficultyWeights  map[models.QuestionDifficulty]float64                                        `yaml:"difficulty_weights"`
	TypeWeights        map[models.QuestionDifficulty]map[models.QuestionType]float64                `yaml:"type_weights"`
	CategoryWeights    map[models.QuestionDifficulty]map[models.QuestionCategory]float64            `yaml:"category_weights"`
	IntentWeights      map[models.QuestionCategory]map[models.QuestionDifficulty]map[string]float64 `yaml:"intent_weights"`
	MCCorrectCountDist map[models.QuestionDifficulty]map[int]float64                                `yaml:"mc_correct_count_dist"`
}

type PoolMapping struct {
	Correct    string `yaml:"correct"`
	Distractor string `yaml:"distractor"`
}

type IntentPoolMapping struct {
	SC *PoolMapping `yaml:"sc,omitempty"`
	MC *PoolMapping `yaml:"mc,omitempty"`
	TF *PoolMapping `yaml:"tf,omitempty"`
}

type PromptTemplateSet struct {
	SingleChoice   []string `yaml:"sc"`
	MultipleChoice []string `yaml:"mc"`
	TrueFalse      []string `yaml:"tf"`
}

func (s PromptTemplateSet) TemplatesFor(qtype models.QuestionType) []string {
	switch qtype {
	case models.QuestionTypeSingleChoice:
		return s.SingleChoice
	case models.QuestionTypeMultipleChoice:
		return s.MultipleChoice
	case models.QuestionTypeTrueFalse:
		return s.TrueFalse
	default:
		return nil
	}
}

type IntentSpec struct {
	OptionKind  string            `yaml:"option_kind"`
	PoolMapping IntentPoolMapping `yaml:"pool_mapping"`
	Templates   PromptTemplateSet `yaml:"templates"`
}

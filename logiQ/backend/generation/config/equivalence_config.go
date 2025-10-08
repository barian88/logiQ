package config

import "backend/models"

type EquivalenceRuleSpec struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type EquivalenceRulesCatalog struct {
	Equivalent    []EquivalenceRuleSpec `yaml:"equivalent"`
	NonEquivalent []EquivalenceRuleSpec `yaml:"nonequivalent"`
}

type EquivalenceDifficultyConfig struct {
	ChainStepsDist map[int]float64 `yaml:"chain_steps_dist"`
}

type EquivalenceConfig struct {
	Rules      EquivalenceRulesCatalog                                   `yaml:"rules"`
	Difficulty map[models.QuestionDifficulty]EquivalenceDifficultyConfig `yaml:"difficulty"`
}

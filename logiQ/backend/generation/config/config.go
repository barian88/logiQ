package config

import (
	"backend/generation/core"
	"backend/models"
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	DifficultyProfiles map[models.QuestionDifficulty]DifficultyProfile `yaml:"difficulty_profiles"`
	Planner            PlannerWeights                                  `yaml:"planner"`
	Intents            map[string]IntentSpec                           `yaml:"intents"`
	Inference          InferenceConfig                                 `yaml:"inference"`
	Equivalence        EquivalenceConfig                               `yaml:"equivalence"`
}

var OpNameToKind = map[string]core.NodeKind{
	"NOT": core.Not,
	"AND": core.And,
	"OR":  core.Or,
	"IMP": core.Impl,
	"IFF": core.Iff,
	"VAR": core.Var, // Added VAR for inf config parsing
}

type plannerConfig struct {
	DifficultyProfiles map[models.QuestionDifficulty]DifficultyProfile `yaml:"difficulty_profiles"`
	Planner            PlannerWeights                                  `yaml:"planner"`
	Intents            map[string]IntentSpec                           `yaml:"intents"`
}

type inferenceFile struct {
	Inference InferenceConfig `yaml:"inference"`
}

type equivalenceFile struct {
	Equivalence EquivalenceConfig `yaml:"equivalence"`
}

func LoadConfig() (AppConfig, error) {
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)

	var cfg AppConfig

	var planner plannerConfig
	if err := loadYAML(filepath.Join(baseDir, "config.yaml"), &planner); err != nil {
		return AppConfig{}, err
	}
	cfg.DifficultyProfiles = planner.DifficultyProfiles
	cfg.Planner = planner.Planner
	cfg.Intents = planner.Intents

	var inf inferenceFile
	if err := loadYAML(filepath.Join(baseDir, "inference.yaml"), &inf); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return AppConfig{}, err
		}
	} else {
		cfg.Inference = inf.Inference
	}

	var eq equivalenceFile
	if err := loadYAML(filepath.Join(baseDir, "equivalence.yaml"), &eq); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return AppConfig{}, err
		}
	} else {
		cfg.Equivalence = eq.Equivalence
	}

	return cfg, nil
}

func loadYAML(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return yaml.Unmarshal(data, out)
}

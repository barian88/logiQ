package sampler

import (
	"backend/generation/config"
	"backend/models"
	"encoding/json"
	"math/rand/v2"
	"testing"
	"time"
)

func TestSample(t *testing.T) {

	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	simpler := NewSampler(cfg)
	rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))

	sampleRes, err := simpler.Sample(rng, models.QuestionDifficultyHard, "", "")
	if err != nil {
		t.Fatalf("failed to sample: %v", err)
	}
	// to json and print plan
	jsonPlan, err := json.Marshal(sampleRes.Plan)
	if err != nil {
		t.Fatalf("failed to marshal plan: %v", err)
	}
	t.Logf("Sampled Plan: %s", jsonPlan)

	// to json and print profile
	jsonProfile, err := json.Marshal(sampleRes.Profile)
	if err != nil {
		t.Fatalf("failed to marshal profile: %v", err)
	}
	t.Logf("Sampled Profile: %s", jsonProfile)

}

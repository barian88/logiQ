package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// QuestionStats holds the statistics for a single question.
type QuestionStats struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	QuestionID     primitive.ObjectID `json:"question_id" bson:"question_id"`
	TotalAnswers   int64              `json:"total_answers" bson:"total_answers"`
	CorrectAnswers int64              `json:"correct_answers" bson:"correct_answers"`
}

type DimensionPortion struct {
	Value   string  `json:"value"`   // e.g., category or difficulty value  like "truthTable" or "easy"
	Count   int64   `json:"count"`   // number of questions in this dimension
	Portion float64 `json:"portion"` // portion of total questions, e.g., 0.25 for 25%
}

// QuestionDistributions holds distributions for multiple question dimensions.
type QuestionDistributions struct {
	Difficulty []DimensionPortion `json:"difficulty"`
	Type       []DimensionPortion `json:"type"`
	Category   []DimensionPortion `json:"category"`
}

// AccuracyByDimension holds accuracy statistics for a single value within a dimension.
type AccuracyByDimension struct {
	Value    string  `json:"value"`    // e.g., "easy"
	Accuracy float64 `json:"accuracy"` // Calculated as CorrectAnswers / TotalAnswers
}

// AccuracyDistributions holds accuracy data across multiple dimensions.
type AccuracyDistributions struct {
	Difficulty []AccuracyByDimension `json:"difficulty"`
	Type       []AccuracyByDimension `json:"type"`
	Category   []AccuracyByDimension `json:"category"`
}

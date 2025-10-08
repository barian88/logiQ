package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type QuestionType string

const (
	QuestionTypeSingleChoice   QuestionType = "singleChoice"
	QuestionTypeMultipleChoice QuestionType = "multipleChoice"
	QuestionTypeTrueFalse      QuestionType = "trueFalse"
)

type QuestionCategory string

const (
	QuestionCategoryTruthTable  QuestionCategory = "truthTable"
	QuestionCategoryEquivalence QuestionCategory = "equivalence"
	QuestionCategoryInference   QuestionCategory = "inference"
)

type QuestionDifficulty string

const (
	QuestionDifficultyEasy   QuestionDifficulty = "easy"
	QuestionDifficultyMedium QuestionDifficulty = "medium"
	QuestionDifficultyHard   QuestionDifficulty = "hard"
)

type Question struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	QuestionText       string             `json:"question_text" bson:"question_text" binding:"required"`
	Options            []string           `json:"options" bson:"options" binding:"required"`
	CorrectAnswerIndex []int              `json:"correct_answer_index" bson:"correct_answer_index" binding:"required"`
	Type               QuestionType       `json:"type" bson:"type" binding:"required,oneof=singleChoice multipleChoice trueFalse"`    // "singleChoice" | "multipleChoice" | "trueFalse"
	Category           QuestionCategory   `json:"category" bson:"category" binding:"required,oneof=truthTable equivalence inference"` // "truthTable" | "equivalence" | "inference"
	Difficulty         QuestionDifficulty `json:"difficulty" bson:"difficulty" binding:"required,oneof=easy medium hard"`             // "easy" | "medium" | "hard"
	IsActive           bool               `json:"is_active" bson:"is_active"`
}

type GenerateQuestionRequest struct {
	Number int `json:"number" bson:"number" binding:"required"`
	// 以下三个字段可选，不提供则表示不限制
	Category   QuestionCategory   `json:"category" bson:"category" binding:"omitempty,oneof=truthTable equivalence inference"`
	Difficulty QuestionDifficulty `json:"difficulty" bson:"difficulty" binding:"omitempty,oneof=easy medium hard"`
	Type       QuestionType       `json:"type" bson:"type" binding:"omitempty,oneof=singleChoice multipleChoice trueFalse"`
}

type GetQuestionListRequest struct {
	Category   QuestionCategory   `json:"category,omitempty" form:"category" bson:"category,omitempty" binding:"omitempty,oneof=truthTable equivalence inference"`
	Difficulty QuestionDifficulty `json:"difficulty,omitempty" form:"difficulty" bson:"difficulty,omitempty" binding:"omitempty,oneof=easy medium hard"`
	Type       QuestionType       `json:"type,omitempty" form:"type" bson:"type,omitempty" binding:"omitempty,oneof=singleChoice multipleChoice trueFalse"`
	Page       int                `json:"page,omitempty" form:"page" bson:"page,omitempty" binding:"omitempty,min=1"`
	PageSize   int                `json:"page_size,omitempty" form:"page_size" bson:"page_size,omitempty" binding:"omitempty,min=1,max=100"`
}

// QuestionResponseForAdmin 用于管理员查看的题目详情，包含正确率
type QuestionResponseForAdmin struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	QuestionText       string             `json:"question_text" bson:"question_text" binding:"required"`
	Options            []string           `json:"options" bson:"options" binding:"required"`
	CorrectAnswerIndex []int              `json:"correct_answer_index" bson:"correct_answer_index" binding:"required"`
	Type               QuestionType       `json:"type" bson:"type" binding:"required,oneof=singleChoice multipleChoice trueFalse"`    // "singleChoice" | "multipleChoice" | "trueFalse"
	Category           QuestionCategory   `json:"category" bson:"category" binding:"required,oneof=truthTable equivalence inference"` // "truthTable" | "equivalence" | "inference"
	Difficulty         QuestionDifficulty `json:"difficulty" bson:"difficulty" binding:"required,oneof=easy medium hard"`             // "easy" | "medium" | "hard"
	IsActive           bool               `json:"is_active" bson:"is_active"`

	// 新增的统计字段
	TotalAnswers   int64   `json:"total_answers" bson:"total_answers"`
	CorrectAnswers int64   `json:"correct_answers" bson:"correct_answers"`
	AccuracyRate   float64 `json:"accuracy_rate" bson:"accuracy_rate"`
}
type BatchDeleteRequest struct {
	IDs []string `json:"ids" binding:"required,dive,hexadecimal"`
}

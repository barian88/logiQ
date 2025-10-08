package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuizQuestion struct {
	Question        *Question `json:"question" bson:"question"`
	UserAnswerIndex []int     `json:"user_answer_index" bson:"user_answer_index"`
	IsCorrect       bool      `json:"is_correct" bson:"is_correct"`
}

type QuizType string

const (
	QuizTypeRandomTasks   QuizType = "randomTasks"
	QuizTypeTopicPractice QuizType = "topicPractice"
	QuizTypeByDifficulty  QuizType = "byDifficulty"
	QuizTypeCustomQuiz    QuizType = "customQuiz"
)

type Quiz struct {
	ID                  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID              primitive.ObjectID `json:"user_id" bson:"user_id"`
	Type                QuizType           `json:"type" bson:"type"` // "randomTasks" | "topicPractice" | "byDifficulty" | "customQuiz"
	Questions           []QuizQuestion     `json:"questions" bson:"questions"`
	CorrectQuestionsNum int                `json:"correct_questions_num" bson:"correct_questions_num"`
	CompletionTime      int                `json:"completion_time" bson:"completion_time"` // 秒
	CompletedAt         time.Time          `json:"completed_at" bson:"completed_at"`
}

type CreateQuizRequest struct {
	Type       QuizType           `json:"type" binding:"required,oneof=randomTasks topicPractice byDifficulty customQuiz"`
	Category   QuestionCategory   `json:"category,omitempty"`   // for topicPractice
	Difficulty QuestionDifficulty `json:"difficulty,omitempty"` // for byDifficulty
}

type SubmitQuizRequest struct {
	Type           QuizType       `json:"type" binding:"required,oneof=randomTasks topicPractice byDifficulty customQuiz"`
	Questions      []QuizQuestion `json:"questions" binding:"required"`
	CompletionTime int            `json:"completion_time" binding:"required"`
}

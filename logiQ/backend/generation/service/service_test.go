package service

import (
	"testing"
)

func TestService(t *testing.T) {
	service := NewService()
	questions, err := service.GenerateQuestion(50, "", "", "")
	if err != nil {
		t.Fatal(" Error generating question: ", err)
	}
	for _, question := range questions {
		if len(questions) == 0 {
			t.Fatal("No questions found")
		}

		t.Logf("Question: %s", question.QuestionText)
		for _, option := range question.Options {
			t.Logf("\nOption: %s", option)
		}
		t.Logf("\nCorrect Answers: %v", question.CorrectAnswerIndex)
		t.Logf("\nType: %s", question.Type)
		t.Logf("\nCategory: %s", question.Category)
		t.Logf("\nDifficulty: %s", question.Difficulty)
	}

}

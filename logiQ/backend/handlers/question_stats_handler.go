package handlers

import (
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QuestionStatsHandler struct {
	questionStatsService *services.QuestionStatsService
}

func NewQuestionStatsHandler(questionStatsService *services.QuestionStatsService) *QuestionStatsHandler {
	return &QuestionStatsHandler{
		questionStatsService: questionStatsService,
	}
}

func (h *QuestionStatsHandler) GetDimensionDistribution(c *gin.Context) {
	distribution, err := h.questionStatsService.CalcQuestionDistributions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, distribution)
}

func (h *QuestionStatsHandler) GetDimensionAccuracy(c *gin.Context) {
	accuracy, err := h.questionStatsService.CalcAccuracyByDimension()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, accuracy)
}

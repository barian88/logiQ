package handlers

import (
	"backend/models"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionHandler struct {
	questionService *services.QuestionService
}

func NewQuestionHandler(questionService *services.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
	}
}

// GenerateQuestion 向题库中增加题目
func (h *QuestionHandler) GenerateQuestion(c *gin.Context) {
	var req models.GenerateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	generatedList, err := h.questionService.GenerateQuestion(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 返回插入的题目数量
	c.JSON(http.StatusCreated, generatedList)
}

// GetQuestionById 根据id获取题目
func (h *QuestionHandler) GetQuestionById(c *gin.Context) {
	questionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	question, err := h.questionService.GetQuestionByID(questionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, question)
}

// GetQuestionList 获取题目列表
func (h *QuestionHandler) GetQuestionList(c *gin.Context) {
	var req models.GetQuestionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questions, total, err := h.questionService.GetQuestionList(req.Category, req.Difficulty, req.Type, req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"list":         questions,
		"total":        total,
		"current_page": req.Page,
	})
}

// BatchDeleteQuestions 批量软删除题目
func (h *QuestionHandler) BatchDeleteQuestions(c *gin.Context) {
	var req models.BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectIDs := make([]primitive.ObjectID, len(req.IDs))
	for i, id := range req.IDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}
		objectIDs[i] = objID
	}

	deletedCount, err := h.questionService.SoftDeleteQuestionsByIDs(objectIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Batch delete successful",
		"deleted_count": deletedCount,
	})
}

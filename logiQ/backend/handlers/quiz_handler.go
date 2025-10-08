package handlers

import (
	"backend/middleware"
	"backend/models"
	"backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type QuizHandler struct {
	quizService *services.QuizService
}

func NewQuizHandler(quizService *services.QuizService) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
	}
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	var req models.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token获取用户ID
	userID, exist := middleware.GetUserIDFromContext(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User authentication information not found"})
		return
	}

	quiz, err := h.quizService.CreateQuiz(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, quiz)
}

func (h *QuizHandler) SubmitQuiz(c *gin.Context) {
	var req models.SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token获取用户ID
	userID, exist := middleware.GetUserIDFromContext(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User authentication information not found"})
		return
	}

	quiz, err := h.quizService.SubmitQuiz(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

func (h *QuizHandler) GetUserQuizHistory(c *gin.Context) {

	// 从JWT token获取用户ID
	userID, exist := middleware.GetUserIDFromContext(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User authentication information not found"})
		return
	}

	quizzes, err := h.quizService.GetUserQuizHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//	返回gin.H{"data": quizzes}格式的JSON响应，包装成一个map，方便前端处理
	c.JSON(http.StatusOK, gin.H{"data": quizzes})
}

func (h *QuizHandler) GetQuiz(c *gin.Context) {
	quizID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	quiz, err := h.quizService.GetQuizByID(quizID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

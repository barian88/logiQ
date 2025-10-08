package handlers

import (
	"backend/middleware"
	"backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserStatsHandler struct {
	userStatsService *services.UserStatsService // 用户统计服务实例，用于处理用户统计相关的业务逻辑
}

// NewUserStatsHandler 创建新的用户统计处理器实例
func NewUserStatsHandler(userStatsService *services.UserStatsService) *UserStatsHandler {
	return &UserStatsHandler{
		userStatsService: userStatsService,
	}
}

// GetUserStats 获取用户统计信息
func (h *UserStatsHandler) GetUserStats(c *gin.Context) {
	// 从JWT token获取用户ID
	userID, exist := middleware.GetUserIDFromContext(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User authentication information not found"})
		return
	}

	stats, err := h.userStatsService.GetUserStatsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Stats not found"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

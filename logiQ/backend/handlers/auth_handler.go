package handlers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器结构体 - 处理所有与认证相关的HTTP请求
type AuthHandler struct {
	userService         *services.UserService         // 用户服务实例，用于处理业务逻辑
	verificationService *services.VerificationService // 验证码服务实例
}

// NewAuthHandler 创建新的认证处理器实例
func NewAuthHandler(userService *services.UserService, verificationService *services.VerificationService) *AuthHandler {
	return &AuthHandler{
		userService:         userService,
		verificationService: verificationService,
	}
}

// RegisterRequest 注册请求接口 - 处理POST /auth/register-request请求
func (h *AuthHandler) RegisterRequest(c *gin.Context) {
	var req models.RegisterRequestModel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data format",
			"error":   err.Error(),
		})
		return
	}

	err := h.userService.RegisterRequest(&req)
	if err != nil {
		var statusCode int
		errorMessage := err.Error()

		if strings.Contains(errorMessage, "already registered") {
			statusCode = http.StatusConflict
		} else {
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"message": "Registration request failed",
			"error":   errorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verification code has been sent to your email",
	})
}

// CompleteRegistration 完成注册接口 - 处理POST /auth/complete-registration请求
func (h *AuthHandler) CompleteRegistration(c *gin.Context) {
	var req models.CompleteRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data format",
			"error":   err.Error(),
		})
		return
	}

	// 解析临时token
	claims, err := utils.ParseTemporaryToken(req.TemporaryToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid or expired temporary token",
			"error":   err.Error(),
		})
		return
	}

	// 验证token用途
	if claims.Purpose != "registration" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid token purpose",
			"error":   "Token is not for registration",
		})
		return
	}

	// 完成注册
	user, err := h.userService.CompleteRegistration(claims.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Registration completion failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login 用户登录接口 - 处理POST /auth/login请求
func (h *AuthHandler) Login(c *gin.Context) {
	// 第1步：解析登录请求数据
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data format",
			"error":   err.Error(),
		})
		return
	}

	// 第2步：验证用户凭据
	user, err := h.userService.Login(&req)
	if err != nil {
		// 登录失败，返回401未授权错误
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Login failed",
			"error":   err.Error(),
		})
		return
	}

	// 第3步：生成JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate access token",
			"error":   err.Error(),
		})
		return
	}

	// 第4步：登录成功，返回用户信息和访问令牌
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// LoginAdmin 管理员登录接口 - 处理POST /auth/login-admin请求
func (h *AuthHandler) LoginAdmin(c *gin.Context) {
	// 第1步：解析登录请求数据
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data format",
			"error":   err.Error(),
		})
		return
	}

	// 第2步：验证用户凭据
	user, err := h.userService.LoginAdmin(&req)
	if err != nil {
		// 登录失败，返回401未授权错误
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Login failed",
			"error":   err.Error(),
		})
		return
	}

	// 第3步：生成JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate access token",
			"error":   err.Error(),
		})
		return
	}

	// 第4步：登录成功，返回访问令牌
	c.JSON(http.StatusOK, token)
}

// Logout 用户登出接口 - 处理POST /auth/logout请求
func (h *AuthHandler) Logout(c *gin.Context) {
	// 对于JWT token，通常在客户端删除token即可
	// 这里可以实现token黑名单等高级功能
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

// SendVerificationCode 发送验证码接口 - 处理POST /auth/send-verification请求
func (h *AuthHandler) SendVerificationCode(c *gin.Context) {
	var req models.SendVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data format",
			"error":   err.Error(),
		})
		return
	}

	err := h.verificationService.SendVerificationCode(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to send verification code",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verification code sent successfully",
	})
}

// VerifyCode 验证验证码接口 - 处理POST /auth/verify-code请求
func (h *AuthHandler) VerifyCode(c *gin.Context) {
	var req models.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data format",
			"error":   err.Error(),
		})
		return
	}

	// 验证purpose参数
	if req.Purpose != "registration" && req.Purpose != "password_reset" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid purpose",
			"error":   "Purpose must be 'registration' or 'password_reset'",
		})
		return
	}

	// 验证验证码
	err := h.verificationService.VerifyCode(req.Email, req.VerificationCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Verification failed",
			"error":   err.Error(),
		})
		return
	}

	// 生成临时token
	temporaryToken, err := utils.GenerateTemporaryToken(req.Email, req.Purpose)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate temporary token",
			"error":   err.Error(),
		})
		return
	}

	// 验证成功，返回临时token
	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"message":         "Verification successful",
		"temporary_token": temporaryToken,
		"expires_in":      300, // 5分钟 = 300秒
	})
}

// UpdatePassword 更新密码接口 - 处理POST /auth/update-password请求
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data format",
			"error":   err.Error(),
		})
		return
	}

	// 解析临时token
	claims, err := utils.ParseTemporaryToken(req.TemporaryToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid or expired temporary token",
			"error":   err.Error(),
		})
		return
	}

	// 验证token用途
	if claims.Purpose != "password_reset" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid token purpose",
			"error":   "Token is not for password reset",
		})
		return
	}

	// 更新密码
	err = h.userService.UpdatePassword(claims.Email, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Password update failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password updated successfully",
	})
}

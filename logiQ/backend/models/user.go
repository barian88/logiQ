package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户数据模型 - 对应MongoDB中的users集合
type User struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`             // MongoDB自动生成的唯一ID
	Username          string             `json:"username" bson:"username"`                       // 用户名
	Email             string             `json:"email" bson:"email"`                             // 用户邮箱，用于登录
	Password          string             `json:"-" bson:"password"`                              // 密码（json:"-"表示不在JSON响应中返回）
	ProfilePictureUrl string             `json:"profile_picture_url" bson:"profile_picture_url"` // 用户头像URL
	Role              string             `json:"role" bson:"role"`                               // 用户角色：user或admin
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`                   // 账户创建时间
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`                   // 最后更新时间
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SendVerificationCodeRequest 发送验证码请求
type SendVerificationCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyCodeRequest 验证验证码请求
type VerifyCodeRequest struct {
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required"`
	Purpose          string `json:"purpose" binding:"required"` // "registration" 或 "password_reset"
}

// CompleteRegistrationRequest 完成注册请求
type CompleteRegistrationRequest struct {
	TemporaryToken string `json:"temporary_token" binding:"required"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	TemporaryToken string `json:"temporary_token" binding:"required"`
	NewPassword    string `json:"new_password" binding:"required,min=6"`
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Username          string `json:"username,omitempty"`
	ProfilePictureUrl string `json:"profile_picture_url,omitempty"`
}

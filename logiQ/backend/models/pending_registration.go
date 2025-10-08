package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PendingRegistration 待注册用户数据模型
type PendingRegistration struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"` // 已加密的密码
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// RegisterRequestModel 注册请求数据模型
type RegisterRequestModel struct {
	Username string `json:"username" binding:"required,min=1"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

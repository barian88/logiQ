package services

import (
	"backend/utils"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"time"
)

// VerificationService 验证码服务
type VerificationService struct {
}

// NewVerificationService 创建验证码服务实例
func NewVerificationService() *VerificationService {
	return &VerificationService{}
}

// SendVerificationCode 发送验证码
func (s *VerificationService) SendVerificationCode(email string) error {
	// 生成4位验证码
	code, err := s.generateCode()
	if err != nil {
		return errors.New("failed to generate verification code")
	}

	// 将验证码存入Redis，有效期5分钟
	// 如果有旧的验证码，此操作会直接覆盖
	err = utils.SetCache(email, code, 5*time.Minute)
	if err != nil {
		return errors.New("failed to save verification code to cache")
	}

	// 发送邮件
	err = utils.SendEmailViaGmail(email, code)
	if err != nil {
		log.Printf("Email sending failed: %v", err)
		return errors.New("failed to send email")
	}

	return nil
}

// VerifyCode 验证验证码
func (s *VerificationService) VerifyCode(email, code string) error {
	// 从Redis获取验证码
	cachedCode, err := utils.GetCache(email)
	if err != nil {
		return errors.New("failed to get verification code from cache")
	}

	// 检查验证码是否存在或已过期
	if cachedCode == "" {
		return errors.New("invalid or expired verification code")
	}

	// 比较验证码
	if cachedCode != code {
		return errors.New("invalid or expired verification code")
	}

	// 验证成功，立即删除验证码
	err = utils.DeleteCache(email)
	if err != nil {
		// 记录日志，但不要因此次失败而导致验证失败，因为密钥无论如何都会过期
		log.Printf("Failed to delete verified code from cache for email %s: %v", email, err)
	}

	return nil
}

// generateCode 生成4位数字验证码
func (s *VerificationService) generateCode() (string, error) {
	const digits = "0123456789"
	result := make([]byte, 4)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		result[i] = digits[num.Int64()]
	}

	return string(result), nil
}

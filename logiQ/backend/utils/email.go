package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendEmailViaGmail 通过Gmail SMTP发送邮件
func SendEmailViaGmail(toEmail, code string) error {
	// Gmail SMTP 配置从环境变量读取
	from := os.Getenv("GMAIL_EMAIL")
	password := os.Getenv("GMAIL_APP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	subject := os.Getenv("EMAIL_SUBJECT")

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" || subject == "" {
		return fmt.Errorf("email configuration is incomplete")
	}

	// 邮件内容
	body := fmt.Sprintf("Your verification code is: %s. Please use it within 5 minutes.", code)
	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body + "\r\n")

	// 认证
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 发送邮件
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, msg)
}

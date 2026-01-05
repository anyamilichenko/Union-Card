package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GeneratePassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}

func SendPasswordEmail(toEmail string, newPassword string) error {
	smtpHost := "smtp.mail.ru"
	smtpPort := "587"
	senderEmail := "anna-100m@mail.ru"
	senderPassword := "eluQAxkn8oDMAUR4F933"

	subject := "Confirmation Code"
	body := fmt.Sprintf("Your password code is: %s", newPassword)
	msg := "From: " + senderEmail + "\n" +
		"To: " + toEmail + "\n" +
		"Subject: " + subject + "\n\n" + body

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{toEmail}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

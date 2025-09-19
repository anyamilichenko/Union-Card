package controllers

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Функция генерации пароля (кода подтверждения)
func GeneratePassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}

// Функция для отправки email
func SendPasswordEmail(toEmail string, newPassword string) error {
	// Настройки SMTP сервера для Mail.ru
	smtpHost := "smtp.mail.ru"
	smtpPort := "587" // Можно использовать порт 465 для SSL
	senderEmail := os.Getenv("LOGIN")
	senderPassword := os.Getenv("PASSWORD")

	// Получатель и тема письма
	subject := "Confirmation Code"
	body := fmt.Sprintf("Your password code is: %s", newPassword)
	msg := "From: " + senderEmail + "\n" +
		"To: " + toEmail + "\n" +
		"Subject: " + subject + "\n\n" + body

	// Аутентификация с использованием ваших данных
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// Отправка письма
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{toEmail}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

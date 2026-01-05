package main

import (
	"bilet/backend"
	"bilet/backend/config"
	"bilet/backend/entity"
	"bilet/backend/handler"
	"bilet/backend/repository"
	"bilet/backend/service"
	"bilet/backend/utils"
	_ "github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация базы данных
	db, err := gorm.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Автомиграция
	db.AutoMigrate(
		&entity.Account{},
		&entity.Token{},
	)

	// Инициализация репозиториев
	accountRepo := repository.NewAccountRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// Запуск очистки токенов
	go tokensCleanup(tokenRepo, cfg.TokensCleanupPeriod)

	// Инициализация утилит
	jwtUtil := utils.NewJWTUtil(cfg.JWTKey, accountRepo, tokenRepo)

	// Инициализация сервисов
	authService := service.NewAuthService(accountRepo, tokenRepo, jwtUtil)
	userService := service.NewUserService(accountRepo)

	// Инициализация хендлеров
	authHandler := handler.NewAuthHandler(authService, jwtUtil)
	userHandler := handler.NewUserHandler(userService, authService)

	// Настройка роутера
	r := backend.SetupRouter(authHandler, userHandler, authService)

	log.Println("Server starting on :1460")
	r.Run(":1460")
}

func tokensCleanup(tokenRepo repository.TokenRepository, period time.Duration) {
	ticker := time.NewTicker(period)
	for range ticker.C {
		tokenRepo.DeleteExpired()
	}
}

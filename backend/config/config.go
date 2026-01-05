package config

import (
	"os"
	"time"
)

const (
	AccessTokenExpiration  = 300 * time.Minute
	RefreshTokenExpiration = 5 * 24 * time.Hour
	TokensCleanupPeriod    = 30 * time.Minute
)

var (
	JwtKey = []byte(os.Getenv("JWT_KEY"))
)

type Config struct {
	DatabaseURL            string
	JWTKey                 []byte
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
	TokensCleanupPeriod    time.Duration
}

func Load() *Config {
	return &Config{
		DatabaseURL:            buildDSN(),
		JWTKey:                 []byte(os.Getenv("JWT_KEY")),
		AccessTokenExpiration:  300 * time.Minute,
		RefreshTokenExpiration: 5 * 24 * time.Hour,
		TokensCleanupPeriod:    30 * time.Minute,
	}
}

func buildDSN() string {
	// Собираем DSN строку из переменных окружения
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "Anna2109"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "ProfBilet"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	return "host=" + host + " user=" + user + " password=" + password +
		" dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=Europe/Moscow"
}

package config

import (
	"os"
	"time"
)

const (
	AccessTokenExpiration                  = 5 * time.Minute
	RefreshTokenExpiration                 = 5 * 24 * time.Hour
	TokensCleanupPeriod                    = 30 * time.Minute
	PasswordCost                           = 14
	PasswordResetCodeLength                = 20
	PasswordResetCodesExpiration           = 15 * time.Minute
	PasswordResetCodesCacheCleanupInterval = 30 * time.Minute
	StateStringLength                      = 20
)

var (
	JwtKey = []byte(os.Getenv("JWT_KEY"))
)

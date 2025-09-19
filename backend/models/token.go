package models

import (
	"time"
)

type Token struct {
	JTI       string `gorm:"primaryKey"`
	Subject   string
	IsRevoked bool `gorm:"default:false"`
	ExpiresAt int64
}

func TokensRepeatedCleanup(period time.Duration) {
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ticker.C:
			DB.Where("expires_at < ?", time.Now().Unix()).Delete(&Token{})
		}
	}
}

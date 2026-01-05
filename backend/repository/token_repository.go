package repository

import (
	"bilet/backend/entity"
	"github.com/jinzhu/gorm"
	"time"
)

type TokenRepository interface {
	Create(token *entity.Token) error
	FindByJTI(jti string) (*entity.Token, error)
	Revoke(jti string) error
	DeleteBySubject(subject string) error
	DeleteExpired()
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(token *entity.Token) error {
	return r.db.Create(token).Error
}

func (r *tokenRepository) FindByJTI(jti string) (*entity.Token, error) {
	var token entity.Token
	err := r.db.Where("jti = ?", jti).First(&token).Error
	return &token, err
}

func (r *tokenRepository) Revoke(jti string) error {
	return r.db.Model(&entity.Token{}).Where("jti = ?", jti).Update("is_revoked", true).Error
}

func (r *tokenRepository) DeleteBySubject(subject string) error {
	return r.db.Where("subject = ?", subject).Delete(&entity.Token{}).Error
}

func (r *tokenRepository) DeleteExpired() {
	r.db.Where("expires_at < ?", time.Now().Unix()).Delete(&entity.Token{})
}

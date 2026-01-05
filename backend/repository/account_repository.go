package repository

import (
	"bilet/backend/entity"
	"errors"
	"github.com/jinzhu/gorm"
)

var ErrNotFound = errors.New("record not found")

type AccountRepository interface {
	Create(account *entity.Account) error
	FindByID(id uint) (*entity.Account, error)
	FindByEmail(email string) (*entity.Account, error)
	Update(account *entity.Account) error
	Delete(id uint) error
	FindAll() ([]entity.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *entity.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) FindByID(id uint) (*entity.Account, error) {
	var account entity.Account
	err := r.db.Where("id = ?", id).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &account, err
}

func (r *accountRepository) FindByEmail(email string) (*entity.Account, error) {
	var account entity.Account
	err := r.db.Where("email = ?", email).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &account, err
}

func (r *accountRepository) Update(account *entity.Account) error {
	return r.db.Save(account).Error
}

func (r *accountRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&entity.Account{}).Error
}

func (r *accountRepository) FindAll() ([]entity.Account, error) {
	var accounts []entity.Account
	err := r.db.Find(&accounts).Error
	return accounts, err
}

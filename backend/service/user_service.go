package service

import (
	"bilet/backend/code"
	"bilet/backend/entity"
	"bilet/backend/repository"
	"errors"
)

type UserService interface {
	GetUserByID(id uint) (*entity.Account, *code.ResultCode)
	GetAllUsers() ([]entity.Account, *code.ResultCode)
	CreateUser(account *entity.Account) (*entity.Account, *code.ResultCode)
	DeleteUser(adminID, targetID uint) *code.ResultCode
	UpdateUser(account *entity.Account) *code.ResultCode
}

type userService struct {
	accountRepo repository.AccountRepository
}

func NewUserService(accountRepo repository.AccountRepository) UserService {
	return &userService{accountRepo: accountRepo}
}

func (s *userService) GetUserByID(id uint) (*entity.Account, *code.ResultCode) {
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return nil, code.InternalServerError.SetMessage(err.Error())
	}
	if account == nil {
		return nil, &code.UserDoesNotExist
	}

	return account, nil
}

func (s *userService) GetAllUsers() ([]entity.Account, *code.ResultCode) {
	accounts, err := s.accountRepo.FindAll()
	if err != nil {
		return nil, code.InternalServerError.SetMessage(err.Error())
	}

	return accounts, nil
}

func (s *userService) CreateUser(account *entity.Account) (*entity.Account, *code.ResultCode) {
	existing, err := s.accountRepo.FindByEmail(account.Email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, code.InternalServerError.SetMessage(err.Error())
	}
	if existing != nil {
		return nil, &code.EmailIsBusy
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, code.InternalServerError.SetMessage(err.Error())
	}

	return account, nil
}

func (s *userService) DeleteUser(adminID, targetID uint) *code.ResultCode {
	// Проверка: админ не может удалить себя
	if adminID == targetID {
		return code.BadRequest.SetMessage("Admin cannot delete themselves")
	}

	admin, err := s.accountRepo.FindByID(adminID)
	if err != nil {
		return code.InternalServerError.SetMessage(err.Error())
	}
	if admin == nil {
		return &code.UserDoesNotExist
	}

	// Проверка: только админы могут удалять пользователей
	if admin.Role != "admin" {
		return &code.Forbidden
	}

	target, err := s.accountRepo.FindByID(targetID)
	if err != nil {
		return code.InternalServerError.SetMessage(err.Error())
	}
	if target == nil {
		return &code.UserDoesNotExist
	}

	if err := s.accountRepo.Delete(targetID); err != nil {
		return code.InternalServerError.SetMessage(err.Error())
	}

	return nil
}

func (s *userService) UpdateUser(account *entity.Account) *code.ResultCode {
	if err := s.accountRepo.Update(account); err != nil {
		return code.InternalServerError.SetMessage(err.Error())
	}

	return nil
}

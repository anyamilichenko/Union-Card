package service

import (
	"bilet/backend/code"
	"bilet/backend/entity"
	"bilet/backend/repository"
	"bilet/backend/utils"
	_ "errors"
	"time"
)

type AuthService interface {
	Login(email, password string) (*entity.Account, *code.ResultCode)
	Logout(token string) *code.ResultCode
	CreateTokens(refreshToken string) (string, string, *code.ResultCode)
	ResetPassword(email string) *code.ResultCode
	ValidateToken(token string) (*entity.Claims, *code.ResultCode)
}

type authService struct {
	accountRepo repository.AccountRepository
	tokenRepo   repository.TokenRepository
	jwtUtil     utils.JWTUtil
}

func NewAuthService(accountRepo repository.AccountRepository, tokenRepo repository.TokenRepository, jwtUtil utils.JWTUtil) AuthService {
	return &authService{
		accountRepo: accountRepo,
		tokenRepo:   tokenRepo,
		jwtUtil:     jwtUtil,
	}
}

func (s *authService) Login(email, password string) (*entity.Account, *code.ResultCode) {
	account, err := s.accountRepo.FindByEmail(email)
	if err != nil {
		return nil, code.InternalServerError.SetMessage(err.Error())
	}
	if account == nil {
		return nil, &code.UserDoesNotExist
	}

	if account.HashedPassword == "" {
		return nil, &code.UserPasswordIsNotSet
	}

	if !utils.IsPasswordCorrect(password, account.HashedPassword) {
		return nil, &code.InvalidPassword
	}

	return account, nil
}

func (s *authService) Logout(token string) *code.ResultCode {
	claims, err := s.jwtUtil.ParseToken(token)
	if err != nil {
		return code.Unauthorized.SetMessage(err.Error())
	}

	// Отзываем токен в базе данных
	if err := s.tokenRepo.Revoke(claims.Id); err != nil {
		return code.InternalServerError.SetMessage("Failed to revoke token")
	}

	return nil
}

func (s *authService) CreateTokens(refreshToken string) (string, string, *code.ResultCode) {
	claims, err := s.jwtUtil.ParseToken(refreshToken)
	if err != nil {
		return "", "", code.Unauthorized.SetMessage(err.Error())
	}

	if time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
		return "", "", code.Unauthorized.SetMessage("Token has expired")
	}

	token, err := s.tokenRepo.FindByJTI(claims.Id)
	if err != nil {
		return "", "", code.Unauthorized.SetMessage("Token not found")
	}

	if token.IsRevoked {
		if err := s.tokenRepo.DeleteBySubject(claims.Subject); err != nil {
			return "", "", code.InternalServerError.SetMessage("Failed to delete user tokens")
		}
		return "", "", code.Unauthorized.SetMessage("Token has already been used")
	}

	newRefreshToken, errCode := s.jwtUtil.NewRefreshToken(claims.Subject, claims.Role, claims.UserId)
	if errCode != nil {
		return "", "", errCode
	}

	newAccessToken, errCode := s.jwtUtil.NewAccessToken(claims.Subject, claims.Role, claims.UserId)
	if errCode != nil {
		return "", "", errCode
	}

	return newRefreshToken, newAccessToken, nil
}

func (s *authService) ResetPassword(email string) *code.ResultCode {
	account, err := s.accountRepo.FindByEmail(email)
	if err != nil {
		return code.InternalServerError.SetMessage(err.Error())
	}
	if account == nil {
		return &code.UserDoesNotExist
	}

	newPassword := utils.GeneratePassword(10)
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return code.InternalServerError.SetMessage("Failed to hash password")
	}

	account.HashedPassword = hashedPassword
	if err := s.accountRepo.Update(account); err != nil {
		return code.InternalServerError.SetMessage("Failed to update password")
	}

	if err := utils.SendPasswordEmail(account.Email, newPassword); err != nil {
		return code.InternalServerError.SetMessage("Failed to send email")
	}

	return nil
}

func (s *authService) ValidateToken(token string) (*entity.Claims, *code.ResultCode) {
	claims, err := s.jwtUtil.ParseToken(token)
	if err != nil {
		return nil, code.Unauthorized.SetMessage(err.Error())
	}

	return claims, nil
}

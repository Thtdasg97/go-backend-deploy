package usecase

import (
	"go-backend/internal/dto"

	"github.com/golang-jwt/jwt/v5"
)

type TokenUsecase interface {
	CreateAccessToken(userId int) (string, error)
	CreateRefreshToken(userId int) (string, error)
	VerifyAccessToken(accessToken string, options ...jwt.ParserOption) (*dto.CustomClaim, error)
	VerifyRefreshToken(refreshToken string, options ...jwt.ParserOption) (*dto.CustomClaim, error)
}

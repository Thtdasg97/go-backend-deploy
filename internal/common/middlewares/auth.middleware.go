package middlewares

import (
	"errors"
	"fmt"
	"go-backend/internal/common/response"
	"go-backend/internal/repository"
	"go-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	tokenUsecase   usecase.TokenUsecase
	userRepository repository.UserRepository
}

func NewAuthMiddleware(tokenUsecase usecase.TokenUsecase, userRepository repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		tokenUsecase:   tokenUsecase,
		userRepository: userRepository,
	}
}

func (a *AuthMiddleware) Protect(ctx *gin.Context) {

	accessToken, err := ctx.Cookie("accessToken")
	if err != nil {
		ctx.Error(response.NewUnauthorizedErrorException())
		ctx.Abort()
		return
	}

	claim, err := a.tokenUsecase.VerifyAccessToken(accessToken)
	if err != nil {

		// trả về 403 khi token hết hạn
		fmt.Println(err)
		if errors.Is(err, jwt.ErrTokenExpired) {
			ctx.Error(response.NewForbiddenErrorException())
			ctx.Abort()
			return
		}

		ctx.Error(response.NewUnauthorizedErrorException())
		ctx.Abort()
		return
	}

	user, err := a.userRepository.FindUserById(ctx, claim.UserId)
	if err != nil {

		ctx.Error(response.NewUnauthorizedErrorException())
		ctx.Abort()
		return
	}

	// Truyền dữ liệu đi dùng .Set
	ctx.Set("user", user)

}

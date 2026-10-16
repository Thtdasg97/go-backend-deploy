package delivery

import (
	"go-backend/internal/common/middlewares"
	"go-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

type authDelivery struct {
	authHandler    *handler.AuthHandler
	authMiddleware *middlewares.AuthMiddleware
}

func NewAuthDelivery(authHandler *handler.AuthHandler, authMiddleware *middlewares.AuthMiddleware) *authDelivery {
	return &authDelivery{
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
	}
}

func (d *authDelivery) RegisterRouter(apiGroup *gin.RouterGroup) {
	authGroup := apiGroup.Group("auth")
	{
		authGroup.POST("register", d.authHandler.Register)
		authGroup.POST("login", d.authHandler.Login)

		// tạo API refresh token dùng khi access token hết hạn sẽ call refresh token
		authGroup.POST("refresh-token", d.authHandler.RefreshToken)
		authGroup.GET("google/login", d.authHandler.GoogleLogin)
		authGroup.GET("google/callback", d.authHandler.GoogleCallback)

		// chỗ này đê thứ tự hơi kỳ hen, ko biết có ý đồ gì không
		authGroup.Use(d.authMiddleware.Protect)
		authGroup.GET("get-info", d.authHandler.GetInfo)
	}
}

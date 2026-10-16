package delivery

import (
	"go-backend/internal/common/middlewares"
	"go-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

type orderDelivery struct {
	orderHandler   *handler.OrderHandler
	authMiddleware *middlewares.AuthMiddleware
}

func NewOrderDelivery(orderHandler *handler.OrderHandler, authMiddleware *middlewares.AuthMiddleware) *orderDelivery {
	return &orderDelivery{
		orderHandler:   orderHandler,
		authMiddleware: authMiddleware,
	}
}

func (d *orderDelivery) RegisterRouter(apiGroup *gin.RouterGroup) {
	orderGroup := apiGroup.Group("order")
	{
		orderGroup.Use(d.authMiddleware.Protect)
		orderGroup.GET("", d.orderHandler.FindAll)
		orderGroup.POST("", d.orderHandler.Create)
	}
}

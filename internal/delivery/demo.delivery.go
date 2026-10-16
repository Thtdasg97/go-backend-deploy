package delivery

import (
	"go-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

type demoDelivery struct {
	demoHandler *handler.DemoHandler
}

func NewDemoDelivery(demoHandler *handler.DemoHandler) *demoDelivery {
	return &demoDelivery{
		demoHandler: demoHandler,
	}
}

func (d *demoDelivery) RegisterRouter(apiGroup *gin.RouterGroup) {
	demoGroup := apiGroup.Group("demo")

	{
		// không truyền hàm vào Query()
		demoGroup.GET("query", d.demoHandler.Query)

		// demo param
		demoGroup.DELETE("param/:id", d.demoHandler.Param)

		// body
		demoGroup.POST("body", d.demoHandler.Body)

		// header
		demoGroup.PATCH("header", d.demoHandler.Header)

	}

}

// demo call back function
func GET(url string, handler func(*gin.Context)) {
	url = "/login"
	// logic

	var context *gin.Context
	handler(context)
}

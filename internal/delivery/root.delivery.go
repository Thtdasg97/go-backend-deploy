package delivery

import (
	"github.com/gin-gonic/gin"
)

type rootDelivery struct {
	demoDelivery        *demoDelivery
	articleDelivery     *articleDelivery
	authDelivery        *authDelivery
	userDelivery        *userDelivery
	chatGroupDelivery   *chatGroupDelivery
	chatMessageDelivery *chatMessageDelivery
	searchDelivery      *searchDelivery
	orderDelivery       *orderDelivery
}

func NewRootDelivery(demoDelivery *demoDelivery,
	articleDelivery *articleDelivery,
	authDelivery *authDelivery, userDelivery *userDelivery,
	chatGroupDelivery *chatGroupDelivery,
	chatMessageDelivery *chatMessageDelivery,
	searchDelivery *searchDelivery,
	orderDelivery *orderDelivery) *rootDelivery {
	return &rootDelivery{
		demoDelivery:        demoDelivery,
		articleDelivery:     articleDelivery,
		authDelivery:        authDelivery,
		userDelivery:        userDelivery,
		chatGroupDelivery:   chatGroupDelivery,
		chatMessageDelivery: chatMessageDelivery,
		searchDelivery:      searchDelivery,
		orderDelivery:       orderDelivery,
	}
}

// hàm đăng ký router
// r là struct (this)
func (r *rootDelivery) RegisterRouter(ginEngine *gin.Engine) {
	apiGroup := ginEngine.Group("api")
	{
		// gom các bộ API
		r.demoDelivery.RegisterRouter(apiGroup)
		r.articleDelivery.RegisterRouter(apiGroup)
		r.authDelivery.RegisterRouter(apiGroup)
		r.userDelivery.RegisterRouter(apiGroup)
		r.chatGroupDelivery.RegisterRouter(apiGroup)
		r.chatMessageDelivery.RegisterRouter(apiGroup)
		r.searchDelivery.RegisterRouter(apiGroup)
		r.orderDelivery.RegisterRouter(apiGroup)

	}

}

//

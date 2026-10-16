package delivery

import (
	"go-backend/internal/common/cache"
	"go-backend/internal/common/middlewares"
	"go-backend/internal/handler"
	"time"

	"github.com/gin-gonic/gin"
)

type articleDelivery struct {
	articleHandler *handler.ArticleHandler
	cache          *cache.Cache
	authMiddleware *middlewares.AuthMiddleware
}

func NewArticleDelivery(articleHandler *handler.ArticleHandler,
	cache *cache.Cache,
	authMiddleware *middlewares.AuthMiddleware) *articleDelivery {
	return &articleDelivery{
		articleHandler: articleHandler,
		cache:          cache,
		authMiddleware: authMiddleware,
	}
}

func (d *articleDelivery) RegisterRouter(apiGroup *gin.RouterGroup) {
	// article là tên của table, sẽ tạo các api tương tác với table này
	articleGroup := apiGroup.Group("article")
	{
		articleGroup.Use(d.authMiddleware.Protect)
		articleGroup.POST("", d.articleHandler.Create)
		articleGroup.GET("", d.cache.CacheRedis(50*time.Second), d.articleHandler.FindAll)
		// có id để user gửi lên mình biết xóa bài nào base trên id
		articleGroup.DELETE(":id", d.articleHandler.Delete)
	}

}

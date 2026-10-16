package dependency

import (
	"go-backend/ent"
	"go-backend/internal/common/cache"
	"go-backend/internal/common/elastic"
	"go-backend/internal/common/env"
	"go-backend/internal/common/middlewares"
	"go-backend/internal/common/rabbitmq"
	"go-backend/internal/common/socket"
	storage_impl "go-backend/internal/common/storage/impl"
	"go-backend/internal/delivery"
	"go-backend/internal/handler"
	"go-backend/internal/repository/repository_impl"
	"go-backend/internal/usecase/usecase_impl"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// define function Injection

func Injection(ginEngine *gin.Engine,
	entClient *ent.Client,
	gormClient *gorm.DB,
	env *env.Env,
	allowOrigins []string,
	rabbitmq *rabbitmq.RabbitMQ) {

	tokenUsecase := usecase_impl.NewTokenUsecase(env)
	cache := cache.NewCache(env)

	// root -> delivery -> handler -> usecase
	//
	articleRepository := repository_impl.NewArticleRepository(entClient, gormClient)
	userRepository := repository_impl.NewUserRepository(entClient)
	chatGroupRepository := repository_impl.NewchatGroupRepository(entClient)
	chatGroupMemberRepository := repository_impl.NewChatGroupMemberRepository(entClient)
	unitOfWorkRepository := repository_impl.NewUnitOfWorkRepository(entClient)
	chatMessageRepository := repository_impl.NewChatMessageRepository(entClient)

	//repository này đặc biệt, chỉ tương tác với elastic
	elastic := elastic.NewElastic(env, articleRepository, userRepository)
	elastic.InitActicle()
	elastic.InitUser()
	searchRepository := repository_impl.NewSearchRepository(elastic)

	//
	localFileStorage := storage_impl.NewLocalFileStorage("public")
	cloudinaryFileStorage := storage_impl.NewCloudinaryStorage(env)

	// xử lý socket
	chatUsecase := usecase_impl.NewChatUsecase(tokenUsecase, userRepository, chatGroupRepository, chatGroupMemberRepository, unitOfWorkRepository, chatMessageRepository)
	chatHanlder := handler.NewChatHandler(chatUsecase)

	socket := socket.NewSocket(chatHanlder)
	socket.Start(ginEngine, allowOrigins)
	authMiddleware := middlewares.NewAuthMiddleware(tokenUsecase, userRepository)

	articleUsecase := usecase_impl.NewArticleUsecase(articleRepository)
	articleHandler := handler.NewArticleHandler(articleUsecase)
	articleDelivery := delivery.NewArticleDelivery(articleHandler, cache, authMiddleware)

	//

	//  new instance usecase
	demoUsecase := usecase_impl.NewDemoUsecase()

	// new instance demo handler
	demoHandler := handler.NewDemoHandler(demoUsecase)

	// new instance demo delivery
	demoDelivery := delivery.NewDemoDelivery(demoHandler)
	// demoDelivery.RegisterRouter(ginEngine)

	authUsecase := usecase_impl.NewAuthUsecase(userRepository, tokenUsecase, env)
	authHandler := handler.NewAuthHandler(authUsecase, env)
	authDelivery := delivery.NewAuthDelivery(authHandler, authMiddleware)

	// bộ user
	userUsecase := usecase_impl.NewUserUsecase(userRepository, localFileStorage, cloudinaryFileStorage)
	userHandler := handler.NewUserHandler(userUsecase)
	// implement bảo vệ API
	userDelivery := delivery.NewUserDelivery(userHandler, authMiddleware)

	// Chat Group
	// chatGroupRepository := repository_impl.NewChatGroupRepository()

	chatGroupUsecase := usecase_impl.NewChatGroupUsecase(chatGroupRepository)
	chatGroupHandler := handler.NewChatGroupHandler(chatGroupUsecase)
	chatGroupDelivery := delivery.NewChatGroupDelivery(chatGroupHandler)

	chatMessageUsecase := usecase_impl.NewChatMessageUsecase(chatMessageRepository)
	chatMessageHandler := handler.NewChatMessageHandler(chatMessageUsecase)
	chatMessageDelivery := delivery.NewChatMessageDelivery(chatMessageHandler)

	// Elastic Search

	searchUsecase := usecase_impl.NewSearchUsecase(searchRepository)
	searchHandler := handler.NewSearchHandler(searchUsecase)
	searchDelivery := delivery.NewSearchDelivery(searchHandler)

	// rabbitMQ - order API

	orderUsecase := usecase_impl.NewOrderUsecase(rabbitmq)
	orderHandler := handler.NewOrderHandler(orderUsecase)
	orderDelivery := delivery.NewOrderDelivery(orderHandler, authMiddleware)

	// new instance lên, rồi call hàm register router
	rootDelivery := delivery.NewRootDelivery(
		demoDelivery,
		articleDelivery,
		authDelivery,
		userDelivery,
		chatGroupDelivery,
		chatMessageDelivery,
		searchDelivery,
		orderDelivery,
	)

	rootDelivery.RegisterRouter(ginEngine)
}

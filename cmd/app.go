package main

import (
	"fmt"
	"go-backend/ent"
	"go-backend/internal/common/ent_client"
	"go-backend/internal/common/env"
	"go-backend/internal/common/gorm_client"
	"go-backend/internal/common/middlewares"
	"go-backend/internal/common/rabbitmq"
	"go-backend/internal/common/response"
	"go-backend/internal/common/swagger"
	dependency "go-backend/internal/di"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Khai báo một strcut tên App
type App struct {
	ginEngine *gin.Engine
	// define thêm biến môi trường trong struct
	env       *env.Env
	entClient *ent.Client
	rabbitmq  *rabbitmq.RabbitMQ
}

// Constructor function - tạo và trả về con trỏ, trỏ đến instance mới của App

func NewApp() *App {

	// New biến env
	env := env.New()

	// tạo instance gin engine
	ginEngine := gin.New()
	// tự add logger và recovery (custom error)
	ginEngine.Use(gin.Logger())
	// gắn middleware
	ginEngine.Use(middlewares.ErrorHandler)
	ginEngine.Use(gin.CustomRecovery(func(ctx *gin.Context, err any) {
		// code theo thiết kế của Gin
		ctx.Error(response.NewInternalServerErrorException())
		ctx.Abort()

	}))

	// relativePath: bí danh dùng ở web, thay thế tên folder
	// root: folder thật ở máy/ server
	ginEngine.Static("images", "./public/images")

	// Implement CORS + hạn chế *
	corsConfig := cors.DefaultConfig()
	// cho phép FE
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	ginEngine.Use(cors.New(corsConfig))

	// Implement CORS = tay
	// viết 1 hàm ẩn danh, set header cho API trả về ,
	// ginEngine.Use(func(ctx *gin.Context) {
	// 	ctx.Header("access-control-allow-origin", "http://localhost:3000")

	// 	if ctx.Request.Method == "OPTIONS" {
	// 		ctx.Header("access-control-allow-headers", "Origin, Content-Length, Content-Type")
	// 		ctx.Header("access-control-allow-methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
	// 		ctx.Header("access-control-max-age", "43200")
	// 		ctx.AbortWithStatus(204)
	// 		return
	// 	}

	// })

	swagger.Start(ginEngine, env)

	// ginEngine.Use(middlewares.A)
	// ginEngine.Use(middlewares.B)
	// ginEngine.Use(middlewares.C)

	// aplly DI

	entClient := ent_client.New(env)
	gormClient := gorm_client.New(env)
	rabbitmq := rabbitmq.NewRabbitMQ(env)
	dependency.Injection(ginEngine, entClient, gormClient, env, corsConfig.AllowOrigins, rabbitmq)

	return &App{
		ginEngine: ginEngine,
		env:       env,
		entClient: entClient,
		rabbitmq:  rabbitmq,
	}
}

// method start được gắn vào receiver *App
// this = struct app
func (a *App) Start() {
	// this.Start()
	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)

	addr := fmt.Sprintf("%s:%s", a.env.Host, a.env.Port)

	a.ginEngine.Run(addr)
}

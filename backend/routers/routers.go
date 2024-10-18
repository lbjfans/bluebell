package routers

import (
	"backend/controller"
	"backend/logger"
	"backend/middlewares"
	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 设置成发布模式，logger写到文件就行；开发者模式写到屏幕和文件，容易debug
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/", func(context *gin.Context) { // router捕获uri，用gin.context处理
		controller.ResponseSuccess(context, "index page!!!")
	})
	// 注册，后序可以加group
	r.POST("/register", controller.SignUpHandler)
	// 登录
	r.POST("/login", controller.LoginHandler)
	r.Use(middlewares.JWTAuthMiddleware()) // 登陆后可以使用的功能
	{
		r.GET("/ping", func(context *gin.Context) {
			controller.ResponseSuccess(context, "pong")
		})
	}
	return r
}

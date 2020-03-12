package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaka53/Short-address/controller"
	"github.com/xiaka53/Short-address/dao"
	"github.com/xiaka53/Short-address/middleware"
)

//路由初始化
func InitRouter(middlewares ...gin.HandlerFunc) (router *gin.Engine) {
	router = gin.Default()
	router.Use(middlewares...)

	app := &controller.App{
		Storage: dao.InitRedis(),
	}
	v1 := router.Group("/api")
	v1.Use(middleware.RecoverMiddleware(), middleware.RequestLog(), middleware.IPAuthMiddleware(), middleware.TranslationMiddleware())
	{
		controller.ApiRegister(app, v1)
	}
	router.GET("/0/:shortLink", app.Rediect)
	return
}

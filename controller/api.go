package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xiaka53/Short-address/dto"
	"github.com/xiaka53/Short-address/intef"
	"github.com/xiaka53/Short-address/middleware"
)

type App struct {
	Storage intef.RedisStorage
}

func ApiRegister(router *gin.RouterGroup) {
	app := App{}
	router.POST("/crater", app.createShortlink)
	router.GET("/info", app.getShortlinkInfo)
}

//创建短链接
func (app *App) createShortlink(c *gin.Context) {
	param := &dto.ShortenReq{}
	if err := param.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, middleware.ParameterError, err)
		return
	}
	newUrl, err := app.Storage.Shorten(param.URL, param.ExpirationInMinuxtes)
	if err != nil {
		middleware.ResponseError(c, 504, err)
		return
	}
	middleware.ResponseSuccess(c, newUrl)
}

//端地址解析
func (app *App) getShortlinkInfo(c *gin.Context) {
	vals := c.Query("shortlink")
	if len(vals) < 5 {
		middleware.ResponseError(c, middleware.ParameterError, errors.New("url error!"))
		return
	}
	urlInfo, err := app.Storage.ShortLinkInfo(vals)
	if err != nil {
		middleware.ResponseError(c, 504, err)
		return
	}
	middleware.ResponseSuccess(c, urlInfo)
}

//重定向
func (app *App) rediect(c *gin.Context) {
	newUrl, err := app.Storage.UnShorten()
	if err != nil {
		middleware.ResponseError(c, 504, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}

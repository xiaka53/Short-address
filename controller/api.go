package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xiaka53/Short-address/dto"
	"github.com/xiaka53/Short-address/intef"
	"github.com/xiaka53/Short-address/middleware"
	"github.com/xiaka53/Short-address/public"
	"net/http"
)

type App struct {
	Storage intef.RedisStorage
}

func ApiRegister(app *App, router *gin.RouterGroup) {
	router.POST("/shorter", app.createShortlink)
	router.GET("/info", app.getShortlinkInfo)
}

//创建短链接
func (app *App) createShortlink(c *gin.Context) {
	param := &dto.ShortenReq{}
	if err := param.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, middleware.ParameterError, err)
		return
	}
	if param.ExpirationInMinuxtes == 0 {
		middleware.ResponseError(c, middleware.ParameterError, errors.New("暂未开放"))
		return
	}
	newUrl, err := app.Storage.Shorten(public.GetGinTraceContext(c), param.URL, param.ExpirationInMinuxtes)
	if err != nil {
		middleware.ResponseError(c, 504, err)
		return
	}
	middleware.ResponseSuccess(c, "http://d.cocofan.cn/"+newUrl)
}

//端地址解析
func (app *App) getShortlinkInfo(c *gin.Context) {
	vals := c.Query("shortlink")
	if len(vals) < 5 {
		middleware.ResponseError(c, middleware.ParameterError, errors.New("url error!"))
		return
	}
	urlInfo, err := app.Storage.ShortLinkInfo(public.GetGinTraceContext(c), vals)
	if err != nil {
		middleware.ResponseError(c, 504, err)
		return
	}
	middleware.ResponseSuccess(c, urlInfo)
}

//重定向
func (app *App) Rediect(c *gin.Context) {
	shortLink := c.Param("shortLink")
	newUrl, err := app.Storage.UnShorten(public.GetGinTraceContext(c), shortLink)
	if err != nil {
		middleware.ResponseError(c, 504, err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, newUrl)
}

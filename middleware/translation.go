package middleware

import (
	"demo/gin-frame/public"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9/translations/en"
	"gopkg.in/go-playground/validator.v9/translations/ja"
	"gopkg.in/go-playground/validator.v9/translations/zh"
)

//设置Translation
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.DefaultQuery("locale", "zh")
		trans, _ := public.Uni.GetTranslator(locale)
		switch locale {
		case "zh":
			zh.RegisterDefaultTranslations(public.Validate, trans)
			break
		case "en":
			en.RegisterDefaultTranslations(public.Validate, trans)
			break
		case "ja":
			ja.RegisterDefaultTranslations(public.Validate, trans)
			break
		default:
			en.RegisterDefaultTranslations(public.Validate, trans)
			break
		}
		//设置trans到context
		c.Set("trans", trans)
		c.Next()
	}
}

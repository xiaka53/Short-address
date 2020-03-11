package public

import (
	"github.com/go-playground/locales"
	en2 "github.com/go-playground/locales/en"
	ja2 "github.com/go-playground/locales/ja"
	zh2 "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

var (
	Uni      *ut.UniversalTranslator
	Validate *validator.Validate
)

//初始化语言包
func InitValidate() {
	var (
		en, zh, ja locales.Translator
	)
	en = en2.New()
	zh = zh2.New()
	ja = ja2.New()
	Uni = ut.New(en, zh, ja)
	Validate = validator.New()
}

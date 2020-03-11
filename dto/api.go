package dto

import (
	"errors"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/xiaka53/Short-address/public"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

type ShortenReq struct {
	URL                  string `json:"url" validate:"required"`
	ExpirationInMinuxtes int64  `json:"expiration_in_minuxtes" validate:"min=0"`
}

func (o *ShortenReq) BindingValidParams(c *gin.Context) (err error) {
	if err = c.ShouldBind(o); err != nil {
		return
	}
	v := c.Value("trans")
	trans, ok := v.(ut.Translator)
	if !ok {
		trans, _ = public.Uni.GetTranslator("zh")
	}
	if err = public.Validate.Struct(o); err != nil {
		errs := err.(validator.ValidationErrors)
		sliceErrs := []string{}
		for _, e := range errs {
			sliceErrs = append(sliceErrs, e.Translate(trans))
		}
		return errors.New(strings.Join(sliceErrs, ","))
	}
	return
}

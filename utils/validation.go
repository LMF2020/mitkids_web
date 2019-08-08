package utils

import (
	"errors"
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	translator "gopkg.in/go-playground/validator.v9/translations/zh"
	"strings"
)

/**
验证api方法并打印错误
*/
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	zh       = zhongwen.New()
)

func ValidateParam(params interface{}) error {

	uni = ut.New(zh, zh)

	trans, _ := uni.GetTranslator("zh")

	validate := validator.New()
	translator.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(params)

	if err != nil {

		errs := err.(validator.ValidationErrors)

		var outerr []string

		for _, err := range errs.Translate(trans) {
			outerr = append(outerr, err)
		}

		Log.Error(outerr)
		return errors.New(strings.Join(outerr, ","))
	}

	return nil
}

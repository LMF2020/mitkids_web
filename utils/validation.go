package utils

import (
	"errors"
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	translator "gopkg.in/go-playground/validator.v9/translations/zh"
	"mitkid_web/utils/log"
	"regexp"
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

		log.Logger.Error(outerr)
		return errors.New(strings.Join(outerr, ","))
	}

	return nil
}

func VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//mobile verify
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

//image verify
func VerifyImageFormat(imageName string) bool {
	regular := `(.*)\.(jpg|gif|jpeg|png)$`

	reg := regexp.MustCompile(regular)
	return reg.MatchString(imageName)
}
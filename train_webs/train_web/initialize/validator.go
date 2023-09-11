package initialize

import (
	"reflect"
	"gotrains/train_webs/train_web/global"
	"strings"

	custom_validator "gotrains/train_webs/train_web/validator"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
)

func InitTrans(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhTrans := zh.New()
		enTrans := en.New()

		uni := ut.New(enTrans, zhTrans, enTrans)
		var ok bool
		global.Trans, ok = uni.GetTranslator(local)
		if !ok {
			return
		}
		switch local {
		case "en":
			en_trans.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			zh_trans.RegisterDefaultTranslations(v, global.Trans)
		default:
			en_trans.RegisterDefaultTranslations(v, global.Trans)
		}
		return
	}
	return
}

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", custom_validator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0}非法手机号码", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
}

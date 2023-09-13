package router

import (
	"gotrains/userpassenger_web/user-web/api"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(r *gin.RouterGroup) {
	br := r.Group("base")
	{
		br.GET("captcha", api.CaptchaGet)
		br.POST("send_msg", api.SendSms)
		br.POST("verify", api.CaptchaVerify)
	}
}

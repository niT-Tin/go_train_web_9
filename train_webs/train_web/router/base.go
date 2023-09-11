package router

import (
	"gotrains/train_webs/train_web/api"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(r *gin.RouterGroup) {
	br := r.Group("base")
	{
		br.GET("captcha", api.CaptchaGet)
		br.POST("send_msg", api.SendSms)
	}
}

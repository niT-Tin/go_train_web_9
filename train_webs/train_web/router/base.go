package router

import (
	"gotrains/train_webs/train_web/api"
	"gotrains/train_webs/train_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitTrainRouter(r *gin.RouterGroup) {
	adt := r.Group("ad")
	{
		adt.GET("trains", middlewares.Cors(), middlewares.Trace(), api.GetTrains)
		adt.POST("adt", middlewares.Cors(), middlewares.Trace(), api.AddTrains)
		// adt.POST("adc", api.AddCarriages)
	}
}

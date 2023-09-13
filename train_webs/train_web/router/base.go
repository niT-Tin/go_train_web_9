package router

import (
	"gotrains/train_webs/train_web/api"

	"github.com/gin-gonic/gin"
)

func InitTrainRouter(r *gin.RouterGroup) {
	adt := r.Group("ad")
	{
		adt.GET("trains", api.GetTrains)
		adt.POST("adt", api.AddTrains)
		// adt.POST("adc", api.AddCarriages)
	}
}

package router

import (
	"gotrains/ticketorder_web/ticketorder-web/api"
	"gotrains/ticketorder_web/ticketorder-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitPassengerRouter(r *gin.RouterGroup) {
	br := r.Group("passenger")
	{
		br.GET("list", middlewares.JWTAuth(), middlewares.Trace(), api.GetPassengerList)
		br.POST("add", middlewares.JWTAuth(), middlewares.Trace(), api.AddPassenger)
		br.DELETE("delete/:id", middlewares.JWTAuth(), middlewares.Trace(), api.DeletePassenger)
	}
}

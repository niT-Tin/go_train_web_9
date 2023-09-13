package router

import (
	"gotrains/ticketorder_web/ticketorder-web/api"
	"gotrains/ticketorder_web/ticketorder-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(g *gin.RouterGroup) {
	orderRouter := g.Group("order")
	{
		orderRouter.POST("create", middlewares.JWTAuth(), middlewares.Trace(), api.CreateOrder)
		orderRouter.GET("list", middlewares.JWTAuth(), middlewares.Trace(), api.GetOrder)
	}
}

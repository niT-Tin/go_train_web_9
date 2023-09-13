package router

import (
	"gotrains/ticketorder_web/ticketorder-web/api"
	"gotrains/ticketorder_web/ticketorder-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(g *gin.RouterGroup) {
	orderRouter := g.Group("order")
	{
		orderRouter.POST("create", middlewares.JWTAuth(), api.CreateOrder)
	}
}

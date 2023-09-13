package router

import (
	"gotrains/train_webs/train_web/api"
	"gotrains/train_webs/train_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitTicketsRouter(g *gin.RouterGroup) {
	trainRouter := g.Group("train")
	{
		trainRouter.GET("tickets", middlewares.JWTAuth(), api.GetTickets)
		trainRouter.GET("seat", middlewares.JWTAuth(), api.GetSeatsByTrain)
	}
}

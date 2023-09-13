package router

import (
	"gotrains/train_webs/train_web/api"
	"gotrains/train_webs/train_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitTicketsRouter(g *gin.RouterGroup) {
	trainRouter := g.Group("train")
	{
		trainRouter.GET("tickets", middlewares.JWTAuth(), middlewares.Trace(), api.GetTickets)
		trainRouter.GET("seat", middlewares.JWTAuth(), middlewares.Trace(), api.GetSeatsByTrain)
		trainRouter.GET("station", middlewares.JWTAuth(), middlewares.Trace(), api.GetStations)
		trainRouter.GET("travel", middlewares.JWTAuth(), middlewares.Trace(), api.Travels)
	}
}

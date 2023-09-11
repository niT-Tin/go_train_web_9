package initialize

import (
	"gotrains/ticketorder_web/ticketorder-web/middlewares"
	mrouter "gotrains/ticketorder_web/ticketorder-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.Cors())
	ApiGroup := router.Group("/u/v1")
	mrouter.InitUserRouter(ApiGroup)
	mrouter.InitBaseRouter(ApiGroup)

	return router
}

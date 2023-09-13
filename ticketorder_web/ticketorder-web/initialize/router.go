package initialize

import (
	"gotrains/ticketorder_web/ticketorder-web/middlewares"
	mrouter "gotrains/ticketorder_web/ticketorder-web/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.Cors())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	ApiGroup := router.Group("/u/v1")
	mrouter.InitPassengerRouter(ApiGroup)
	mrouter.InitOrderRouter(ApiGroup)

	return router
}

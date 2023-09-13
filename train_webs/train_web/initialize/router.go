package initialize

import (
	"gotrains/train_webs/train_web/middlewares"
	mrouter "gotrains/train_webs/train_web/router"
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
	mrouter.InitTicketsRouter(ApiGroup)
	mrouter.InitTrainRouter(ApiGroup)

	return router
}

package initialize

import (
	"gotrains/userpassenger_web/user-web/middlewares"
	mrouter "gotrains/userpassenger_web/user-web/router"
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
	mrouter.InitUserRouter(ApiGroup)
	mrouter.InitBaseRouter(ApiGroup)

	return router
}

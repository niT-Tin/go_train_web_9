package initialize

import (
	"gotrains/userpassenger_web/user-web/middlewares"
	mrouter "gotrains/userpassenger_web/user-web/router"

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

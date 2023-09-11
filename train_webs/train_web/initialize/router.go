package initialize

import (
	"gotrains/train_webs/train_web/middlewares"
	mrouter "gotrains/train_webs/train_web/router"

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

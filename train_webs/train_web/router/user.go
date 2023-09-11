package router

import (
	"gotrains/train_webs/train_web/api"
	"gotrains/train_webs/train_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(g *gin.RouterGroup) {
	UserRouter := g.Group("user")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PasswordLogin)
		UserRouter.POST("register", api.Register)
	}
}

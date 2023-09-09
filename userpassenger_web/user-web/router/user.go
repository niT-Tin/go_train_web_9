package router

import (
	"gotrains/userpassenger_web/user-web/api"
	"gotrains/userpassenger_web/user-web/middlewares"

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

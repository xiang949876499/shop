package router

import (
	"github.com/gin-gonic/gin"

	"shop/user_web/api"
	"shop/user_web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	{
		userRouter.GET("list", middlewares.JWTAuth(), api.GetUserList)
		userRouter.POST("pwd_login", api.PassWordLogin)
		userRouter.POST("register", api.Register)

	}
}

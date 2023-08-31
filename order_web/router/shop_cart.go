package router

import (
	"shop/order_web/api/shop_cart"
	"shop/order_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	cartRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth()).Use(middlewares.Trace())
	{
		cartRouter.GET("", shop_cart.List)
		cartRouter.DELETE("/:id", shop_cart.Delete)
		cartRouter.POST("", shop_cart.New)
		cartRouter.PATCH("/:id", shop_cart.Update)
	}
}

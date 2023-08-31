package router

import (
	"shop/goods_web/api/goods"
	"shop/goods_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.List)                                 //商品列表
		GoodsRouter.POST("", middlewares.JWTAuth(), goods.New)          //改接口需要管理员权限
		GoodsRouter.GET("/:id", goods.Detail)                           //获取商品的详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), goods.Delete) //删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)                    //获取商品的库存

		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), goods.Update)
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), goods.UpdateStatus)
	}
}

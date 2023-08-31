package router

import (
	"shop/goods_web/api/banners"
	"shop/goods_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners").Use()
	{
		BannerRouter.GET("", banners.List)                                 // 轮播图列表页
		BannerRouter.DELETE("/:id", middlewares.JWTAuth(), banners.Delete) // 删除轮播图
		BannerRouter.POST("", middlewares.JWTAuth(), banners.New)          //新建轮播图
		BannerRouter.PUT("/:id", middlewares.JWTAuth(), banners.Update)    //修改轮播图信息
	}
}

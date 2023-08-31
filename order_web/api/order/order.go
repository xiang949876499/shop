package order

import (
	"context"
	"net/http"
	"shop/order_web/forms"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	"shop/order_web/api"
	"shop/order_web/global"
	"shop/order_web/models"
	"shop/order_web/proto"
)

func List(ctx *gin.Context) {
	userid, _ := ctx.Get("userid")
	claims, _ := ctx.Get("claims")
	model := claims.(models.CustomClaims)
	request := proto.OrderFilterRequest{}
	if model.AuthorityId == 1 {
		request.UserId = int32(userid.(uint))
	}

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	request.Pages = int32(pagesInt)
	request.PagePerNums = int32(perNumsInt)

	res, err := global.OrderSrvClient.OrderList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": res.Total,
	}
	orderList := make([]interface{}, 0)
	for _, item := range res.Data {
		tmpMap := map[string]interface{}{}
		tmpMap["id"] = item.Id
		tmpMap["status"] = item.Status
		tmpMap["pay_type"] = item.PayType
		tmpMap["user"] = item.UserId
		tmpMap["post"] = item.Post
		tmpMap["total"] = item.Total
		tmpMap["address"] = item.Address
		tmpMap["name"] = item.Name
		tmpMap["mobile"] = item.Mobile
		tmpMap["order_sn"] = item.OrderSn
		tmpMap["id"] = item.Id
		tmpMap["add_time"] = item.AddTime
		orderList = append(orderList, tmpMap)
	}
	reMap["data"] = orderList
	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	orderForm := forms.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}
	userId, _ := ctx.Get("userId")
	res, err := global.OrderSrvClient.CreateOrder(context.Background(),
		&proto.OrderRequest{
			UserId:  int32(userId.(uint)),
			Name:    orderForm.Name,
			Mobile:  orderForm.Mobile,
			Address: orderForm.Address,
			Post:    orderForm.Post,
		})
	if err != nil {
		zap.S().Errorw("新建订单失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         res.Id,
		"alipay_url": "支付地址",
	})
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	//如果是管理员用户则返回所有的订单
	request := proto.OrderRequest{
		Id: int32(i),
	}
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	res, err := global.OrderSrvClient.OrderDetail(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	reMap := gin.H{}
	reMap["id"] = res.OrderInfo.Id
	reMap["status"] = res.OrderInfo.Status
	reMap["user"] = res.OrderInfo.UserId
	reMap["post"] = res.OrderInfo.Post
	reMap["total"] = res.OrderInfo.Total
	reMap["address"] = res.OrderInfo.Address
	reMap["name"] = res.OrderInfo.Name
	reMap["mobile"] = res.OrderInfo.Mobile
	reMap["pay_type"] = res.OrderInfo.PayType
	reMap["order_sn"] = res.OrderInfo.OrderSn

	goodsList := make([]interface{}, 0)
	for _, item := range res.Goods {
		tmpMap := gin.H{
			"id":    item.GoodsId,
			"name":  item.GoodsName,
			"image": item.GoodsImage,
			"price": item.GoodsPrice,
			"nums":  item.Nums,
		}

		goodsList = append(goodsList, tmpMap)
	}
	reMap["goods"] = goodsList
	ctx.JSON(http.StatusOK, reMap)
}

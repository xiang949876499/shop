package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"google.golang.org/grpc/codes"

	"shop/order_srv/global"
	"shop/order_srv/model"
	protobuf "shop/order_srv/proto"
)

type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	Ctx         context.Context
}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	json.Unmarshal(msg.Body, &orderInfo)

	var goodsIds []int32              //查询物品信息
	goodsNum := make(map[int32]int32) //选中商品的数量
	var carts []model.ShoppingCart
	var orderGoods []*model.OrderGoods
	var goodsInfo []*protobuf.GoodsInvInfo

	if res := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}); res.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "没有选中结算的商品"
		return primitive.RollbackMessageState
	}

	for _, cart := range carts {
		goodsIds = append(goodsIds, cart.Goods)
		goodsNum[cart.Goods] = cart.Nums
	}

	//获取商品信息
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &protobuf.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "批量查询商品失败"
		return primitive.RollbackMessageState
	}
	var amount float32
	for _, v := range goods.Data {
		amount += v.ShopPrice * float32(goodsNum[v.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      v.Id,
			GoodsName:  v.Name,
			GoodsImage: v.GoodsFrontImage,
			GoodsPrice: v.ShopPrice,
			Nums:       goodsNum[v.Id],
		})

		goodsInfo = append(goodsInfo, &protobuf.GoodsInvInfo{
			GoodsId: v.Id,
			Num:     goodsNum[v.Id],
		})
	}

	//调用库存信息
	if _, err := global.InventorySrvClient.Sell(context.Background(), &protobuf.SellInfo{GoodsInfo: goodsInfo}); err != nil {
		o.Code = codes.ResourceExhausted
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}

	tx := global.DB.Begin()
	orderInfo.OrderMount = amount
	o.OrderAmount = amount
	o.ID = orderInfo.ID
	if res := global.DB.Save(&orderInfo); res.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "订单创建失败"
		return primitive.RollbackMessageState
	}

	for _, good := range orderGoods {
		good.Order = orderInfo.ID
	}
	if res := global.DB.CreateInBatches(orderGoods, 100); res.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "订单创建失败"
		return primitive.RollbackMessageState
	}

	if res := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(model.ShoppingCart{}); res.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "删除购物车记录失败"
		return primitive.RollbackMessageState
	}
	tx.Commit()
	o.Code = codes.OK
	return primitive.CommitMessageState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	json.Unmarshal(msg.Body, &orderInfo)

	if result := global.DB.Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
		return primitive.CommitMessageState //你并不能说明这里就是库存已经扣减了
	}

	return primitive.RollbackMessageState
}

func InitRocketmq() {
	adds := fmt.Sprintf("%s:%d", global.ServerConfig.Rocketmq.Host, global.ServerConfig.Rocketmq.Port)
	p, err := rocketmq.NewTransactionProducer(&OrderListener{}, producer.WithNameServer([]string{adds}))
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		panic(err)
	}
	global.Rocketmq = p
}

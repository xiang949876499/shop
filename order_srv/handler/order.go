package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop/order_srv/global"
	"shop/order_srv/model"
	protobuf "shop/order_srv/proto"
)

type OrderServer struct {
	protobuf.UnimplementedOrderServer
}

func GenerateOrderSn(userId int32) string {
	//订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

// CartItemList 获取用户的购物车列表
func (*OrderServer) CartItemList(ctx context.Context, req *protobuf.UserInfo) (*protobuf.CartItemListResponse, error) {
	var carts []model.ShoppingCart
	var rsp protobuf.CartItemListResponse

	if res := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&carts); res.Error != nil {
		return nil, res.Error
	} else {
		rsp.Total = int32(res.RowsAffected)
	}

	for _, cart := range carts {
		rsp.Data = append(rsp.Data, &protobuf.ShopCartInfoResponse{
			Id:      cart.ID,
			UserId:  cart.User,
			GoodsId: cart.Goods,
			Nums:    cart.Nums,
			Checked: cart.Checked,
		})
	}

	return &rsp, nil
}

// CreateCartItem 将商品添加到购物车
func (*OrderServer) CreateCartItem(ctx context.Context, req *protobuf.CartItemRequest) (*protobuf.ShopCartInfoResponse, error) {
	var cart model.ShoppingCart
	if res := global.DB.Where(&model.ShoppingCart{User: req.UserId, Goods: req.GoodsId}).First(&cart); res.RowsAffected == 1 {
		cart.Nums += req.Nums
	} else {
		cart.User = req.UserId
		cart.Goods = req.GoodsId
		cart.Nums = req.Nums
		cart.Checked = false
	}

	global.DB.Save(&cart)

	return &protobuf.ShopCartInfoResponse{Id: cart.ID}, nil
}

// UpdateCartItem 更新购物车信息
func (*OrderServer) UpdateCartItem(ctx context.Context, req *protobuf.CartItemRequest) (*emptypb.Empty, error) {
	var cart model.ShoppingCart
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&cart); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	cart.Checked = req.Checked
	if req.Nums > 0 {
		cart.Nums = req.Nums
	}

	global.DB.Save(&cart)
	return &emptypb.Empty{}, nil
}

// DeleteCartItem 删除购物车
func (*OrderServer) DeleteCartItem(ctx context.Context, req *protobuf.CartItemRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).
		Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}

// CreateOrder 创建订单
func (*OrderServer) CreateOrder(ctx context.Context, req *protobuf.OrderRequest) (*protobuf.OrderInfoResponse, error) {
	//var goodsIds []int32              //查询物品信息
	//goodsNum := make(map[int32]int32) //选中商品的数量
	//var carts []model.ShoppingCart
	//var orderGoods []*model.OrderGoods
	//var goodsInfo []*protobuf.GoodsInvInfo
	//
	//if res := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}); res.RowsAffected == 0 {
	//	return nil, status.Errorf(codes.InvalidArgument, "没有选中结算的商品")
	//}
	//
	//for _, cart := range carts {
	//	goodsIds = append(goodsIds, cart.Goods)
	//	goodsNum[cart.Goods] = cart.Nums
	//}
	//
	////获取商品信息
	//goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &protobuf.BatchGoodsIdInfo{Id: goodsIds})
	//if err != nil {
	//	return nil, status.Errorf(codes.Internal, "批量查询商品失败")
	//}
	//var amount float32
	//for _, v := range goods.Data {
	//	amount += v.ShopPrice * float32(goodsNum[v.Id])
	//	orderGoods = append(orderGoods, &model.OrderGoods{
	//		Goods:      v.Id,
	//		GoodsName:  v.Name,
	//		GoodsImage: v.GoodsFrontImage,
	//		GoodsPrice: v.ShopPrice,
	//		Nums:       goodsNum[v.Id],
	//	})
	//
	//	goodsInfo = append(goodsInfo, &protobuf.GoodsInvInfo{
	//		GoodsId: v.Id,
	//		Num:     goodsNum[v.Id],
	//	})
	//}
	//
	////调用库存信息
	//if _, err := global.InventorySrvClient.Sell(context.Background(), &protobuf.SellInfo{GoodsInfo: goodsInfo}); err != nil {
	//	return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败")
	//}
	//
	//tx := global.DB.Begin()
	//
	//orderInfo := model.OrderInfo{
	//	OrderSn:      GenerateOrderSn(req.UserId),
	//	User:         req.UserId,
	//	Address:      req.Address,
	//	SignerName:   req.Name,
	//	SingerMobile: req.Mobile,
	//	Post:         req.Post,
	//}
	//
	//if res := global.DB.Save(&orderInfo); res.RowsAffected == 0 {
	//	tx.Rollback()
	//}
	//for _, good := range orderGoods {
	//	good.Order = orderInfo.ID
	//}
	//if res := global.DB.CreateInBatches(orderGoods, 100); res.RowsAffected == 0 {
	//	tx.Rollback()
	//}
	//
	//if res := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(model.ShoppingCart{}); res.RowsAffected == 0 {
	//	tx.Rollback()
	//}

	orderListener := OrderListener{Ctx: ctx}
	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}
	//应该在消息中具体指明一个订单的具体的商品的扣减情况
	jsonString, _ := json.Marshal(order)

	_, err := global.Rocketmq.SendMessageInTransaction(context.Background(),
		primitive.NewMessage("order_reback", jsonString))
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
		return nil, status.Error(codes.Internal, "发送消息失败")
	}
	if orderListener.Code != codes.OK {
		return nil, status.Error(orderListener.Code, orderListener.Detail)
	}

	return &protobuf.OrderInfoResponse{
		Id:      orderListener.ID,
		OrderSn: order.OrderSn,
		Total:   orderListener.OrderAmount,
	}, nil
}

// OrderList 获取用户的订单列表
func (*OrderServer) OrderList(ctx context.Context, req *protobuf.OrderFilterRequest) (*protobuf.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp protobuf.OrderListResponse

	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)

	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)

	for _, order := range orders {
		rsp.Data = append(rsp.Data, &protobuf.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2023-01-01 15:04:05"),
		})
	}

	return &rsp, nil
}

// OrderDetail 获取订单详情
func (*OrderServer) OrderDetail(ctx context.Context, req *protobuf.OrderRequest) (*protobuf.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp protobuf.OrderInfoDetailResponse

	if res := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id},
		User: req.UserId}).First(&order); res.Error != nil {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	rsp.OrderInfo.Id = order.ID
	rsp.OrderInfo.UserId = order.User
	rsp.OrderInfo.OrderSn = order.OrderSn
	rsp.OrderInfo.PayType = order.PayType
	rsp.OrderInfo.Status = order.Status
	rsp.OrderInfo.Post = order.Post
	rsp.OrderInfo.Total = order.OrderMount
	rsp.OrderInfo.Address = order.Address
	rsp.OrderInfo.Name = order.SignerName
	rsp.OrderInfo.Mobile = order.SingerMobile
	rsp.OrderInfo.AddTime = order.CreatedAt.Format("2023-01-01 15:04:05")

	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}

	for _, good := range orderGoods {
		rsp.Goods = append(rsp.Goods, &protobuf.OrderItemResponse{
			Id:         good.ID,
			OrderId:    good.Order,
			GoodsId:    good.Goods,
			GoodsName:  good.GoodsName,
			GoodsImage: good.GoodsImage,
			GoodsPrice: good.GoodsPrice,
			Nums:       good.Nums,
		})
	}

	return &rsp, nil
}

// UpdateOrderStatus 更新订单状态
func (*OrderServer) UpdateOrderStatus(ctx context.Context, req *protobuf.OrderStatus) (*emptypb.Empty, error) {
	if res := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).
		Update("status", req.Status); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &emptypb.Empty{}, nil
}

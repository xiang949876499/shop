package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop/inventory_srv/global"
	"shop/inventory_srv/model"
	protobuf "shop/inventory_srv/proto"
)

type InventoryServer struct {
	protobuf.UnimplementedInventoryServer
}

// SetInv 设置库存
func (*InventoryServer) SetInv(ctx context.Context, req *protobuf.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	global.DB.First(&inv, req.GoodsId)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num
	global.DB.Save(&inv)

	return &emptypb.Empty{}, nil
}

// InvDetail 查询库存
func (*InventoryServer) InvDetail(ctx context.Context, req *protobuf.GoodsInvInfo) (*protobuf.GoodsInvInfo, error) {
	var inv model.Inventory

	if res := global.DB.First(&inv, req.GoodsId); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}

	return &protobuf.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks - inv.Freeze,
	}, nil
}

// Sell 扣减库存
func (*InventoryServer) Sell(ctx context.Context, req *protobuf.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()
	sell := model.StockSellDetail{
		OrderSn: req.OrderSn,
		Status:  1,
	}
	var details []model.GoodsDetail
	for _, item := range req.GoodsInfo {
		var inv model.Inventory
		details = append(details, model.GoodsDetail{
			Goods: item.GoodsId,
			Num:   item.Num,
		})

		mutex := global.Rs.NewMutex(fmt.Sprintf("goods_%d", item.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}
		if result := global.DB.Where(&model.Inventory{Goods: item.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		if inv.Stocks < item.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		inv.Stocks -= item.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	sell.Detail = details
	if res := tx.Create(&sell); res.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("保存库存扣减历史失败 订单号: %s", req.OrderSn))
	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

// Reback 回归库存
func (*InventoryServer) Reback(ctx context.Context, req *protobuf.SellInfo) (*emptypb.Empty, error) {
	for _, item := range req.GoodsInfo {
		var inv model.Inventory
		if res := global.DB.First(&inv, item.GoodsId); res.RowsAffected == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		inv.Stocks += item.Num
		global.DB.Save(&inv)
	}
	return &emptypb.Empty{}, nil
}

func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	for i := range msgs {
		var orderInfo OrderInfo
		if err := json.Unmarshal(msgs[i].Body, &orderInfo); err != nil {
			global.Log.Errorf("解析json失败 %v \n", msgs[i].Body)
			return consumer.ConsumeSuccess, nil
		}

		tx := global.DB.Begin()
		var sell model.StockSellDetail
		if res := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn,
			Status: 1}).First(&sell); res.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}

		for _, detail := range sell.Detail {
			if res := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods: detail.Goods}).Update("stocks",
				gorm.Expr("stocks+?", detail.Num)); res.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}
		sell.Status = 2
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn}).Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}

/*********2pc做法***********/

// TrySell 预扣减库存
func (*InventoryServer) TrySell(ctx context.Context, req *protobuf.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()
	for _, item := range req.GoodsInfo {
		var inv model.Inventory

		mutex := global.Rs.NewMutex(fmt.Sprintf("goods_%d", item.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}
		if result := global.DB.Where(&model.Inventory{Goods: item.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		if inv.Stocks < item.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//inv.Stocks -= item.Num
		inv.Freeze += item.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

// ConfirmSell 实际扣减库存
func (*InventoryServer) ConfirmSell(ctx context.Context, req *protobuf.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()
	for _, item := range req.GoodsInfo {
		var inv model.Inventory

		mutex := global.Rs.NewMutex(fmt.Sprintf("goods_%d", item.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}
		if result := global.DB.Where(&model.Inventory{Goods: item.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		if inv.Stocks < item.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		inv.Stocks -= item.Num
		inv.Freeze -= item.Num
		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

// CancelSell 撤销库存扣减
func (*InventoryServer) CancelSell(ctx context.Context, req *protobuf.SellInfo) (*emptypb.Empty, error) {
	for _, item := range req.GoodsInfo {
		var inv model.Inventory
		if res := global.DB.First(&inv, item.GoodsId); res.RowsAffected == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		inv.Freeze -= item.Num
		global.DB.Save(&inv)
	}
	return &emptypb.Empty{}, nil
}

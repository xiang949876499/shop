package model

import (
	"database/sql/driver"
	"encoding/json"
)

// todo 还可以设计异地仓库的处理  就近发货

type GoodsDetail struct {
	Goods int32
	Num   int32
}
type GoodsDetailList []GoodsDetail

func (g GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"` //comment 商品编号
	Stocks  int32 `gorm:"type:int"`       // 库存
	Version int32 `gorm:"type:int"`       //分布式锁的乐观锁
	Freeze  int32 `gorm:"type:int"`       //冻结库存
}

type StockSellDetail struct {
	OrderSn string          `gorm:"type:varchar(200);index:idx_order_sn,unique;"`
	Status  int32           `gorm:"type:int"` //1 表示已扣减 2. 表示已归还
	Detail  GoodsDetailList `gorm:"type:varchar(200)"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}

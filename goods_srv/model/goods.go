package model

// Category 分类
type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null" json:"name"`
	ParentCategoryID int32       `json:"parent"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;not null;default:1" json:"level"`
	IsTab            bool        `gorm:"default:false;not null" json:"is_tab"`
}

// Brands 品牌
type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}

// GoodsCategoryBrand 商品分类/品牌对应关系
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands   Brands
}

// TableName 重载表名
func (GoodsCategoryBrand) TableName() string {
	return "bloodstained"
}

// Banner 轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

// Goods 商品详情
type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null comment 是否上架"`
	ShipFree bool `gorm:"default:false;not null	comment 是否免运费"`
	IsNew    bool `gorm:"default:false;not null	comment 是否新品"`
	IsHot    bool `gorm:"default:false;not null	comment 是否热卖"`

	Name            string   `gorm:"type:varchar(50);not null"`
	GoodsSn         string   `gorm:"type:varchar(50);not null comment 内部编号"`
	ClickNum        int32    `gorm:"type:int;default:0;not null comment 点击数"`
	SoldNum         int32    `gorm:"type:int;default:0;not null comment 销量"`
	FavNum          int32    `gorm:"type:int;default:0;not null comment 收藏量"`
	MarketPrice     float32  `gorm:"not null comment 原价"`
	ShopPrice       float32  `gorm:"not null comment 实际价格"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null comment 简介"`
	Images          GormList `gorm:"type:varchar(1000);not null	"`
	DescImages      GormList `gorm:"type:varchar(1000);not null "`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null comment 封面图"`
}

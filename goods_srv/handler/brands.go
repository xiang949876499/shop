package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"goods_srv/global"
	"goods_srv/model"
	protobuf "goods_srv/proto"
)

// BrandList 品牌和轮播图
func (s *GoodsServer) BrandList(c context.Context, req *protobuf.BrandFilterRequest) (*protobuf.BrandListResponse, error) {
	brandListResponse := protobuf.BrandListResponse{}
	var brands []model.Brands
	res := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)

	if res.Error != nil {
		return nil, res.Error
	}

	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)
	brandListResponse.Total = int32(total)

	var brandRes []*protobuf.BrandInfoResponse
	for _, brand := range brands {
		brandRes = append(brandRes, &protobuf.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListResponse.Data = brandRes
	return &brandListResponse, nil
}
func (s *GoodsServer) CreateBrand(c context.Context, req *protobuf.BrandRequest) (*protobuf.BrandInfoResponse, error) {
	if result := global.DB.Where("name=?", req.Name).First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	brand := &model.Brands{}
	brand.Name = req.Name
	brand.Logo = req.Logo
	global.DB.Save(&brand)

	return &protobuf.BrandInfoResponse{Id: brand.ID}, nil
}
func (s *GoodsServer) DeleteBrand(c context.Context, req *protobuf.BrandRequest) (*emptypb.Empty, error) {
	if res := global.DB.Delete(&model.Brands{}, req.Id); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &emptypb.Empty{}, nil
}
func (s *GoodsServer) UpdateBrand(c context.Context, req *protobuf.BrandRequest) (*emptypb.Empty, error) {
	brand := &model.Brands{}
	if res := global.DB.First(brand, req.Id); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	if req.Name != "" {
		brand.Name = req.Name
	}
	if req.Logo != "" {
		brand.Logo = req.Logo
	}

	global.DB.Save(brand)
	return &emptypb.Empty{}, nil
}

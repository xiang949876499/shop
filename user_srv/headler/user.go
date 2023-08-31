package headler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"shop/user_srv/global"
	"shop/user_srv/model"
	"shop/user_srv/proto"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToRsponse(user model.User) proto.UserInfoResponse {
	//在grpc的message中字段有默认值，你不能随便赋值nil进去，容易出错
	//这里要搞清， 哪些字段是有默认值
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	res := global.DB.Find(&users)
	if res.Error != nil {
		return nil, res.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(res.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userinfo := ModelToRsponse(user)
		rsp.Data = append(rsp.Data, &userinfo)
	}

	return rsp, nil
}
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	res := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	if res.Error != nil {
		return nil, res.Error
	}
	rsp := ModelToRsponse(user)
	return &rsp, nil
}
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	res := global.DB.First(&user, req.Id)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	if res.Error != nil {
		return nil, res.Error
	}
	rsp := ModelToRsponse(user)
	return &rsp, nil
}
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	res := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if res.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName

	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	res = global.DB.Create(&user)
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	res := global.DB.First(&user, req.Id)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.AlreadyExists, "用户不存在")
	}
	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender

	res = global.DB.Save(user)
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	passWordInfo := strings.Split(req.EncryptedPassword, "$")
	options := &password.Options{16, 100, 32, sha512.New}
	check := password.Verify(req.Password, passWordInfo[2], passWordInfo[3], options)
	return &proto.CheckResponse{Success: check}, nil
}

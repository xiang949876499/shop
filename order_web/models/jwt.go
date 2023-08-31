package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	BufferTime  int64 //更新token的时间
	jwt.RegisteredClaims
}

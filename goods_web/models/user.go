package models

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJson() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2023-02-17"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32    `json:"id"`
	NickName string   `json:"name"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}

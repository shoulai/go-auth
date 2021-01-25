package model

import (
	"encoding/json"
	"time"
)

type Session struct {
	Id string `json:"id,omitempty"`
	//登录设备
	Device string `json:"device,omitempty"`
	//当前登录用户
	User *User `json:"user,omitempty"`
	//临时身份
	TempUser *User `json:"tempUser,omitempty"`
	//session创建时间
	CreateTime time.Time `json:"createTime,omitempty"`
	//存储额外数据
	Data map[string]interface{} `json:"data,omitempty"`
}

func (v Session) MarshalBinary() (data []byte, err error) {
	return json.Marshal(v)
}

func (v *Session) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, v)
	return err
}

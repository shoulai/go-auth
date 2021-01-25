package model

import (
	"encoding/json"
)

type IUser interface {
	GetId() string

	GetUsername() string

	GetNickname() string

	GetRoles() *[]IRole

	GetMenus() interface{}
}

type User struct {
	Id       string      `json:"id,omitempty"`
	Username string      `json:"username,omitempty"`
	Nickname string      `json:"nickname",omitempty`
	Roles    []Role      `json:"roles,omitempty"`
	Menus    interface{} `json:"menus",omitempty`
}

func (v User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(v)
}

func (v *User) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, v)
	return err
}

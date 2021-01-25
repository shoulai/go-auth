package test_model

import (
	"errors"
	"github.com/shoulai/go-auth/model"
)

var Users = []TestUser{
	{
		Id: "1", Username: "admin", Password: "123456", Nickname: "超级管理员",
		Roles: []TestRole{{Id: "1", Name: "admin", Nickname: "超级管理员",
			Permissions: []TestPermission{
				{Permission: "get:test"},
				{Permission: "get:hello"},
				{Permission: "get:user"},
				{Permission: "post:login"},
				{Permission: "get:session"},
			},
		},
		},
	},
	{
		Id: "2", Username: "user1", Password: "123456", Nickname: "普通用户1",
		Roles: []TestRole{
			{Id: "2", Name: "role_user1", Nickname: "用户1",
				Permissions: []TestPermission{
					{Permission: "get:test"},
					{Permission: "get:user"},
				},
			}},
	},
	{Id: "3", Username: "user2", Password: "123456", Nickname: "普通用户2",
		Roles: []TestRole{
			{Id: "3", Name: "role_user2", Nickname: "用户2",
				Permissions: []TestPermission{
					{Permission: "get:hell"},
					{Permission: "get:user"},
				},
			},
		},
	},
}

func GetUser(username string) (*TestUser, error) {
	for _, user := range Users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, errors.New("用户不存在")
}
func GetUserById(id string) (*TestUser, error) {
	for _, user := range Users {
		if user.Id == id {
			return &user, nil
		}
	}
	return nil, errors.New("用户不存在")
}

type TestUser struct {
	Id       string
	Device   string `form:"device"binding:"required"`
	Username string `form:"username"binding:"required"`
	Password string `form:"password"binding:"required"`
	Nickname string
	Roles    []TestRole
}

func (u TestUser) GetId() string {
	return u.Id
}

func (u TestUser) GetUsername() string {
	return u.Username
}

func (u TestUser) GetNickname() string {
	return u.Nickname
}

func (u TestUser) GetRoles() *[]model.IRole {
	var roles []model.IRole
	for _, val := range u.Roles {
		roles = append(roles, val)
	}
	return &roles
}
func (u TestUser) GetMenus() interface{} {
	return &[]Menu{
		{
			Name:     "首页",
			Sort:     1,
			Router:   "/",
			Icon:     "index",
			ParentId: "",
		},
		{
			Name:     "关于我",
			Sort:     2,
			Router:   "/about",
			Icon:     "about",
			ParentId: "",
		},
	}
}

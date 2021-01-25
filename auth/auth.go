package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/shoulai/go-auth/model"
)

type IAuth interface {

	//登录，返回结果:token,错误信息
	Login(c *gin.Context, user model.IUser, device string) (string, error)

	//临时身份登录
	TempLogin(c *gin.Context, user model.IUser) error

	//登出
	Logout(c *gin.Context) error

	//踢人下线
	ForcedLogout(token string) error

	//判断是否登录
	IsLogin(c *gin.Context) (bool, error)

	//用户信息
	GetCurrentUser(c *gin.Context) (*model.User, error)

	//获取Session
	GetSession(c *gin.Context) (model.Session, error)

	//从session中取数据
	GetSessionData(c *gin.Context, key string) (interface{}, error)

	//向session中存储数据
	SetSessionData(c *gin.Context, key string, value interface{}) error

	//判断用户是否有权限
	Permission(c *gin.Context) bool

	//判断当前访问路径是否支持匿名访问
	Anonymous(c *gin.Context) bool
}

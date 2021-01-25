package impl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shoulai/go-auth/auth"
	"github.com/shoulai/go-auth/cache"
	"github.com/shoulai/go-auth/config"
	"github.com/shoulai/go-auth/model"
	"github.com/shoulai/go-auth/str"
	"github.com/shoulai/go-auth/util/json"
	"github.com/shoulai/go-auth/util/uuid"
	"strings"
	"time"
)

type Auth struct {
	Store       cache.IStore
	Config      *config.Config
	Permissions []model.Permission
}

func New(cfg *config.Config) (auth.IAuth, func(), error) {
	//初始化数据库
	store, f, err := cache.NewStore(*cfg)
	if err != nil {
		return nil, nil, err
	}
	a := &Auth{Store: store, Config: cfg, Permissions: cfg.Permissions}
	return a, f, nil
}

func (a Auth) Login(c *gin.Context, user model.IUser, device string) (string, error) {
	//如果有该用户登录信息，则清除

	//获取当前用户session
	var tokens str.Strings
	var newTokens str.Strings
	err := a.Store.Get(user.GetId(), &tokens)
	if err == nil {
		for _, token := range tokens {
			var session model.Session
			err := a.Store.Get(token, &session)
			if err == nil {
				if a.Config.ConcurrentLogin { //开启并发登录
					newTokens = append(newTokens, token)
				} else { //关闭并发登录
					if session.Device == device {
						a.delSession(token)
					} else {
						newTokens = append(newTokens, token)
					}
				}
			}
		}
	}

	//为当前登录用户创建session
	//生成token
	token := uuid.MustUUID().String()
	//设置Cookie
	a.setCookie(c, token)
	c.Set(a.Config.TokenName, token)
	//构建用户角色列表
	newTokens = append(newTokens, token)
	//构建用户session
	session := model.Session{
		Id:         token,
		CreateTime: time.Now(),
		User: &model.User{
			Id:       user.GetId(),
			Username: user.GetUsername(),
			Nickname: user.GetNickname(),
			Roles:    a.getRoles(user),
			Menus:    user.GetMenus(),
		},
		Device: device,
	}
	//保存当前创建的session到缓存
	//根据sessionId，存储session
	err = a.Store.Set(session.Id, session, a.Config.TokenTimeOut)
	if err != nil {
		return token, err
	}

	//根据用户ID，存储sessionids
	err = a.Store.Set(session.User.Id, newTokens, a.Config.TokenTimeOut)
	if err != nil {
		return token, err
	}
	return token, nil
}

func (a Auth) getRoles(user model.IUser) []model.Role {
	IRoles := user.GetRoles()
	var roles []model.Role
	for _, role := range *IRoles {
		roles = append(roles, model.Role{Id: role.GetId(), Name: role.GetName(), Nickname: role.GetNickname(), Permissions: role.GetPermissions()})
	}
	return roles
}

func (a Auth) setCookie(c *gin.Context, token string) {
	c.SetCookie(a.Config.TokenName, token, a.Config.TokenTimeOut, a.Config.Cookie.Path, a.Config.Cookie.Domain, false, a.Config.Cookie.HttpOnly)
}

func (a Auth) TempLogin(c *gin.Context, user model.IUser) error {
	session, err := a.GetSession(c)

	if err != nil {
		return errors.New("切换临时身份失败")
	}
	//传过来的是当前用户
	if user.GetId() == session.User.Id {
		return errors.New("不支持切换此用户")
	}

	//构建临时用户信息
	session.TempUser = &model.User{
		Id:       user.GetId(),
		Username: user.GetUsername(),
		Nickname: user.GetNickname(),
		Roles:    a.getRoles(user),
		Menus:    user.GetMenus(),
	}

	return a.Store.Set(session.Id, session, a.Config.TokenTimeOut)
}

func (a Auth) Logout(c *gin.Context) error {
	session, err := a.GetSession(c)
	if err != nil {
		return err
	}

	//当前使用的临时身份，退出，先退出临时身份。
	if session.TempUser != nil {
		session.TempUser = nil
		a.Store.Set(session.Id, session, a.Config.TokenTimeOut)
		return nil
	}
	//未使用临时身份，直接退出当前登录
	if session.TempUser == nil {
		a.Store.Del(session.Id)
		var tokens str.Strings
		err := a.Store.Get(session.User.Id, &tokens)
		if err != nil {
			return err
		}
		var newTokens []string
		for _, token := range tokens {
			if token != session.Id {
				newTokens = append(newTokens, token)
			}
		}
		if len(newTokens) < 1 {
			a.Store.Del(session.User.Id)
		} else {
			a.Store.Set(session.User.Id, newTokens, a.Config.TokenTimeOut)
		}
	}
	return err
}

func (a Auth) ForcedLogout(token string) error {
	//获取session
	var session model.Session
	a.Store.Get(token, &session)

	//获取sessionIds
	var strings str.Strings
	err := a.Store.Get(session.User.Id, &strings)
	if err != nil {
		return err
	}

	//从sessionIds中移除登出的sessionId
	var newTokens str.Strings
	for _, token := range strings {
		newTokens = append(newTokens, token)
	}
	//有多个用户在登录
	if len(newTokens) > 0 {
		a.Store.Set(session.User.Id, newTokens, a.Config.TokenTimeOut)
		a.Store.Del(token)
	} else {
		//只有一个用户在登录
		a.Store.Del(session.User.Id)
		a.Store.Del(token)
	}
	return nil
}

func (a Auth) IsLogin(c *gin.Context) (bool, error) {

	//从请求中获取token
	token, err := a.getToken(c)
	if err != nil {
		return false, err
	}

	//从缓存中获取session
	var session model.Session
	err = a.Store.Get(token, &session)
	if err != nil {
		return false, err
	}

	//更新cookie过期时间
	a.setCookie(c, token)

	//更新session过期时间
	session.CreateTime = time.Now()
	err = a.Store.Set(token, session, a.Config.TokenTimeOut)
	if err != nil {
		return true, err
	}

	//从缓存中获取sessionIds
	var tokens str.Strings
	err = a.Store.Get(session.User.Id, &tokens)
	if err != nil {
		return true, err
	}

	//更新存储sessionIds 过期时间
	err = a.Store.Set(session.User.Id, tokens, a.Config.TokenTimeOut)
	if err != nil {
		return true, err
	}
	return true, nil
}

func (a Auth) GetCurrentUser(c *gin.Context) (*model.User, error) {
	session, err := a.GetSession(c)
	if err != nil {
		return nil, errors.New("获取当前用户失败")
	}
	if session.TempUser == nil {
		return session.User, nil
	}
	return session.TempUser, nil
}

func (a Auth) GetSession(c *gin.Context) (model.Session, error) {
	var session model.Session
	//获取token
	token, err := a.getToken(c)
	if err != nil {
		return session, err
	}

	//根据token获取session
	err = a.Store.Get(token, &session)
	return session, err
}

//从session中取数据
func (a Auth) GetSessionData(c *gin.Context, key string) (interface{}, error) {

	//获取session
	session, err := a.GetSession(c)
	if err != nil {
		return nil, err
	}

	//获取session中存储的数据
	value, ok := session.Data[key]
	if !ok {
		return nil, errors.New("not find data")
	}
	return value, err
}

//向session中添加数据
func (a Auth) SetSessionData(c *gin.Context, key string, value interface{}) error {
	//获取token
	token, err := a.getToken(c)
	if err != nil {
		return err
	}

	//根据token获取session
	var session model.Session
	err = a.Store.Get(token, &session)
	if err != nil {
		return err
	}

	//将数据存储进session
	if session.Data == nil {
		m := map[string]interface{}{key: value}
		session.Data = m
	} else {
		session.Data[key] = value
	}

	//将缓存中session更新及刷新过期时间
	err = a.Store.Set(token, session, a.Config.TokenTimeOut)
	if err != nil {
		return err
	}

	//根据用户id获取tokens
	var tokens str.Strings
	err = a.Store.Get(session.User.Id, &tokens)
	if err != nil {
		return err
	}

	//刷新tokens过期时间
	err = a.Store.Set(session.User.Id, tokens, a.Config.TokenTimeOut)
	if err != nil {
		return err
	}

	return nil
}

func (a Auth) Permission(c *gin.Context) bool {
	session, err := a.GetSession(c)
	if err != nil {
		return false
	}
	url := c.Request.URL.Path
	method := c.Request.Method
	//设置当前访问位置，需要的权限
	var perm string
	for _, permission := range a.Permissions {
		if permission.Url == url && strings.ToUpper(method) == strings.ToUpper(permission.Method) {
			//支持匿名访问
			if permission.Anonymous {
				return true
			}
			perm = permission.Permission
			break
		}
	}

	//当前位置尚未设置权限，返回未授权
	if perm == "" {
		return false
	}

	//临时身份登录，则使用临时身份鉴权
	if session.TempUser != nil {
		user := session.TempUser
		for _, role := range user.Roles {
			for _, permission := range role.Permissions {
				if perm == permission {
					return true
				}
			}
		}
	}

	//不是临时用户，则用当前登录用户鉴权
	if session.TempUser == nil {
		user := session.User
		for _, role := range user.Roles {
			for _, permission := range role.Permissions {
				if perm == permission {
					return true
				}
			}
		}
	}
	//当前用户和临时用户，都无该权限，则直接返回未授权
	return false
}

//匿名访问
func (a Auth) Anonymous(c *gin.Context) bool {
	//请求路径
	url := c.Request.URL.Path
	//请求方式
	method := c.Request.Method
	for _, permission := range a.Permissions {
		if strings.ToUpper(permission.Url) == strings.ToUpper(url) &&
			strings.ToUpper(method) == strings.ToUpper(permission.Method) {
			//支持匿名访问
			if permission.Anonymous {
				return true
			}
		}
	}
	return false
}

func (a Auth) setUserid(key string, val interface{}, TimeOut int) error {
	toJson, err := json.MarshalToJson(val)
	if err != nil {
		return err
	}
	return a.Store.Set(key, toJson, TimeOut)
}

//从gin.Context中获取token
func (a Auth) getToken(c *gin.Context) (string, error) {
	token := c.GetString(a.Config.TokenName)
	if token != "" {
		return token, nil
	}
	token, err := c.Cookie(a.Config.TokenName)
	if err == nil {
		return token, nil
	}
	token = c.GetHeader(a.Config.TokenName)
	if token != "" {
		return token, nil
	}
	return token, errors.New("token or cookie invalid")
}

func (a Auth) delTokens(userId string) error {
	return a.Store.Del(userId)
}

func (a Auth) delSession(token string) error {
	return a.Store.Del(token)
}

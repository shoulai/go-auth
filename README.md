# Golang 认证框架(基于gin)

### 尚处于测试使用阶段

```
安装: go get github.com/shoulai/go-auth
```

# 一、项目结构

```
├── auth：认证和鉴权
├── config: 配置中心
├── model：用户,角色，权限，session
├── cache：缓存机制
│   ├── memory：使用sync.Map
│   └── redis：使用redis
├── util: 工具   
├──example 案例
```

# 二、功能介绍：

注意:先实例化auth

```golang
auth, f, err := impl.New(&config.Config{
		CacheName:       "redis",
		TokenName:       "SessionId",
		TokenTimeOut:    5000,
        //开启并发登录
		ConcurrentLogin: true,
		//系统权限
		Permissions: []model.Permission{
			{Url: "/v1/temp/login", Method: "post", Permission: "post:login"},
			{Url: "/v1/session", Method: "get", Permission: "get:session"},
			{Url: "/v1/user", Method: "get", Permission: "get:user"},
			{Url: "/v1/test", Method: "get", Permission: "get:test"},
			{Url: "/v1/hello", Method: "get", Permission: "get:hello"},
		},
		Cookie: &config.Cookie{
			Path:     "/",
			Domain:   "",
			HttpOnly: true,
		},
		//如果需要用redis，将里面参数填写完成，并将Enable设置成true即可
		Redis: &config.Redis{
			Addr: "192.168.101.254:16379",
			DB:   0,
		},
	})
```

## 1.登录

```golang
  auth.Login(c *gin.Context, user model.IUser, device string) error
```
    
## 2.登出
```golang
  auth.Logout(c *gin.Context) error
```
## 3.认证

```golang
 auth.IsLogin(c *gin.Context) bool
```

## 4.鉴权
```golang
 auth.Permission(c *gin.Context) bool
```

## 5.临时身份

```golang
 auth.TempLogin(c *gin.Context, user model.IUser) error
```

## 6.并发登录控制

```golang
  config.ConcurrentLogin: true,
```

## 7.踢人下线

```golang
 auth.ForcedLogout(token string) error
```

# 三、使用案例
```golang
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shoulai/go-auth/config"
	"github.com/shoulai/go-auth/core/model"
	"github.com/shoulai/go-auth/core/subject"
	"github.com/shoulai/go-auth/example/middleware"
	"github.com/shoulai/go-auth/example/test_model"
	"log"
	"net/http"
)

func main() {
	engine := gin.New()
	//初始化认证管理
	auth, _, err := subject.New(&config.Config{
		//系统权限
		Permissions: []model.Permission{
			{Url: "/v1/temp/login", Method: "post", Permission: "post:login"},
			{Url: "/v1/session", Method: "get", Permission: "get:session"},
			{Url: "/v1/user", Method: "get", Permission: "get:user"},
			{Url: "/v1/test", Method: "get", Permission: "get:test"},
			{Url: "/v1/hello", Method: "get", Permission: "get:hello"},
		},
		TimeOut: 1200,
		Cookie: &config.Cookie{
			Name:     "SessionId",
			Path:     "/",
			Domain:   "",
			HttpOnly: true,
		},
		//如果需要用redis，将里面参数填写完成，并将Enable设置成true即可
		Redis: &config.Redis{
			Enable: false,
			Addr:   "192.168.101.254:16379",
			DB:     0,
		},
	})
	if err != nil {
		log.Fatalf("异常:{%v}", err.Error())
	}
	//认证中间件
	engine.Use(middleware.Auth(auth))

	//登录
	engine.POST("/login", func(c *gin.Context) {
		var user test_model.TestUser
		c.ShouldBindQuery(&user)
		queryUser, err2 := test_model.GetUser(user.Username)
		if err2 != nil {
			c.JSON(http.StatusOK, gin.H{"data": gin.H{"code": 201, "message": err2.Error()}})
			c.Abort()
			return
		}
		if user.Password != queryUser.Password {
			c.JSON(http.StatusOK, gin.H{"data": gin.H{"code": 201, "message": "密码错误"}})
			c.Abort()
			return
		}
		err := auth.Login(c, queryUser, user.Device)
		if err != nil {
			log.Printf(err.Error())
		}
		err = auth.SetSessionData(c, "hello", "hello world")
		if err != nil {
			log.Printf(err.Error())
		}
		err = auth.SetSessionData(c, "test", "test world")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"data": gin.H{"code": 201, "message": err.Error()}})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": gin.H{"code": 200, "message": "登录成功"}})
		}
	})

	//退出登录
	engine.GET("/logout", func(c *gin.Context) {
		auth.Logout(c)
		c.JSON(http.StatusOK, gin.H{"data": gin.H{"code": 200, "message": "登出成功"}})
	})

	//登录临时用户
	engine.POST("/v1/temp/login", func(c *gin.Context) {
		value := c.Query("userId")
		if value == "" {
			c.JSON(400, gin.H{"data": gin.H{"code": 400, "message": "userId:不能为空"}})
			return
		}
		user, err2 := test_model.GetUserById(value)
		if err2 != nil {
			c.JSON(200, gin.H{"data": gin.H{"code": 201, "message": err2.Error()}})
			return
		}
		auth.TempLogin(c, user)
		c.JSON(200, gin.H{"data": gin.H{"code": 200, "message": "切换身份成功"}})
	})

	//获取session
	engine.GET("/v1/session", func(c *gin.Context) {
		session, _ := auth.GetSession(c)
		c.JSON(200, gin.H{"data": gin.H{"code": 200, "message": session}})
	})

	//获取当前用户
	engine.GET("/v1/user", func(c *gin.Context) {
		auth, _ := auth.GetCurrentUser(c)
		c.JSON(http.StatusOK, auth)
	})
	//获取hello
	engine.GET("/v1/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"data": gin.H{"code": 200, "message": "hello"}})
	})
	//获取test
	engine.GET("/v1/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"data": gin.H{"code": 200, "message": "test"}})
	})

	engine.Run(":9090")
}

```
具体参考example包下
    


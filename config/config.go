package config

import "github.com/shoulai/go-auth/model"

type Config struct {
	*Redis
	*Cookie
	TokenTimeOut    int
	ConcurrentLogin bool
	CacheName       string
	TokenName       string
	Permissions     []model.Permission
}

type Cookie struct {
	Path     string
	Domain   string
	HttpOnly bool
}

// Config redis配置参数
type Redis struct {
	Enable    bool
	Addr      string // 地址(IP:Port)
	DB        int    // 数据库
	Password  string // 密码
	KeyPrefix string // 存储key的前缀
}

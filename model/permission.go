package model

type Permission struct {
	//请求路径
	Url string
	//请求方式
	Method string
	//请求权限
	Permission string
	//支持匿名访问
	Anonymous bool
}

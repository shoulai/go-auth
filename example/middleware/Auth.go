package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/shoulai/go-auth/auth"
	"net/http"
)

func Auth(auth auth.IAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetPath(c) == "/login" {
			c.Next()
			return
		}

		//匿名访问
		if auth.Anonymous(c) {
			c.Next()
			return
		}

		//需要认证才可访问
		if isLogin, _ := auth.IsLogin(c); !isLogin {
			c.JSON(http.StatusUnauthorized, gin.H{"data": gin.H{"code": 401, "message": "请先登录"}})
			c.Abort()
			return
		}

		//登录后可登出
		if GetPath(c) == "/logout" {
			c.Next()
			return
		}

		//拥有对应权限才可访问
		if isPermission := auth.Permission(c); !isPermission {
			c.JSON(http.StatusForbidden, gin.H{"data": gin.H{"code": 403, "message": "权限不足"}})
			c.Abort()
			return
		}
	}
}

func GetPath(c *gin.Context) string {
	return c.Request.URL.Path
}

package jwt

import (
	"demo/2.gin-demo/pkg/e"
	log "demo/2.gin-demo/pkg/logging"
	"demo/2.gin-demo/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/18 16:24 
    @File: jwt.go    
*/

//新增的权限处理中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				log.Error("Token鉴权失败...err = [%v]", err)
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				log.Warn("Token超时...")
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}
		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}
		c.Next()
	}
}

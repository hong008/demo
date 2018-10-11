package app

import (
	"demo/2.gin-demo/pkg/e"
	"github.com/gin-gonic/gin"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/26 18:17 
    @File: response.go    
*/

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  e.GetMsg(errCode),
		"data": data,
	})
}

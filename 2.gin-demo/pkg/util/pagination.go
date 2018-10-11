package util

import (
	"demo/2.gin-demo/pkg/setting"
	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/18 01:00 
    @File: pagination.go    
*/

//项目page作统一处理
func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}
	return result
}

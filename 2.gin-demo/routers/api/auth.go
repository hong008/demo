package api

import (
	"demo/2.gin-demo/models"
	"demo/2.gin-demo/pkg/e"
	log "demo/2.gin-demo/pkg/logging"
	"demo/2.gin-demo/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/18 16:31 
    @File: auth.go    
*/

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary 请求Token
// @Produce  json
// @Param username query string true "Username"
// @Param password query int false "Password"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS

	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				log.Error("获取Token失败...[%v]", err)
				code = e.ERROR_AUTH_TOKEN
			} else {
				log.Info("Token获取成功...Token = [%v]", token)
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			log.Error("GetAuth err = [%v]", err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

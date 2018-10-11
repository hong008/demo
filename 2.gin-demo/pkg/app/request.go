package app

import (
	"demo/2.gin-demo/pkg/logging"
	"github.com/astaxie/beego/validation"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/26 18:06 
    @File: request.go    
*/

func MarkErrors(erros []*validation.Error) {
	for _, err := range erros {
		logging.Info(err.Key, err.Message)
	}
}
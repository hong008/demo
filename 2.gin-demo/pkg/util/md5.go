package util

import (
	"crypto/md5"
	"encoding/hex"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/26 15:09 
    @File: md5.go    
*/

//MD5加密
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

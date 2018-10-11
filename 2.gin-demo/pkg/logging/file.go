package logging

import (
	"demo/2.gin-demo/pkg/file"
	"fmt"
	"os"
	"time"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/18 16:52 
    @File: file.go    
*/

var (
	RuntimeRootPath string
	LogSavePath     string
	LogSaveName     string
	LogFileExt      string
	TimeFormat      string
)

//获取日志路径
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", RuntimeRootPath, LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s.%s", //年月日.log
		time.Now().Format(TimeFormat),
		LogFileExt,
	)
}

//打开日志文件
func openLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}
	src := dir + "/" + filePath
	perm := file.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %v", src)
	}
	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}
	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to OpenFile: %v", err)
	}
	return f, nil

}

package logging

import (
	"demo/2.gin-demo/pkg/setting"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/18 16:52
    @File: log.go
*/

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func Setup() {
	RuntimeRootPath = setting.AppSetting.RuntimeRootPath
	LogSavePath = setting.AppSetting.LogSavePath
	LogSaveName = setting.AppSetting.LogSaveName
	LogFileExt = setting.AppSetting.LogFileExt
	TimeFormat = setting.AppSetting.TimeFormat

	//日志文件的路径
	filePath := getLogFilePath()
	//日志文件名
	fileName := getLogFileName()

	var err error
	F, err = openLogFile(fileName, filePath)
	if err != nil {
		log.Fatalln(err)
	}
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v...)
}

func Warn(v ...interface{}) {
	setPrefix(WARN)
	logger.Println(v...)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v...)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Println(v...)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}

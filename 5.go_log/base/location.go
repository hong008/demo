package base

import (
	"runtime"
	"strings"
)

//获取调用代码的位置
func GetInvokeerLocation(skipNumber int) (funcPath string, fileName string, line int) {
	pc, file, line, ok := runtime.Caller(skipNumber)
	if !ok {
		return "", "", -1
	}
	if index := strings.LastIndex(file, "/"); index > 0 {
		fileName = file[index+1:]
	}
	funcPtr := runtime.FuncForPC(pc)
	if funcPtr != nil {
		funcPath = funcPtr.Name()
	}
	return funcPath, fileName, line
}

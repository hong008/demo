package lib

import "time"

type RawReq struct {
	ID  int64
	Req []byte
}

type RawResp struct {
	ID     int64
	Resp   []byte
	Err    error
	Elapse time.Duration
}

type RetCode int

const (
	RET_CODE_SUCCESS              RetCode = 0    //成功
	RET_CODE_WARNING_CALL_TIMEOUT         = 1001 //超时警告
	RET_CODE_ERROR_CALL                   = 2001 //调用错误
	RET_CODE_ERROR_RESPONSE               = 2002 //响应内容错误
	RET_CODE_ERROR_CALEE                  = 2003 //被调用方(被测软件)的内部错误
	RET_CODE_FATAL_CALL                   = 3001 //调用过程中发生的致命错误
)

func GetRetCodePlain(code RetCode) string {
	var codePlain string
	switch code {
	case RET_CODE_SUCCESS:
		codePlain = "Success"
	case RET_CODE_WARNING_CALL_TIMEOUT:
		codePlain = "Call Timeout Warning"
	case RET_CODE_ERROR_CALL:
		codePlain = "Call Error"
	case RET_CODE_ERROR_RESPONSE:
		codePlain = "Response Error"
	case RET_CODE_ERROR_CALEE:
		codePlain = "Callee Error"
	case RET_CODE_FATAL_CALL:
		codePlain = "Call Fatal Error"
	default:
		codePlain = "Unknown Result Code"
	}
	return codePlain
}

type CallResult struct {
	ID     int64
	Req    RawReq
	Resp   RawResp
	Code   RetCode
	Msg    string
	Elapse time.Duration
}

const (
	STATUS_ORIGINAL uint32 = 0 //初始状态
	STATUS_STARTING uint32 = 1 //正在启动
	STATUS_STARTED  uint32 = 2 //已启动
	STATUS_STOPPING uint32 = 3 //正在停止
	STATUS_STOPPED  uint32 = 4 //已停止
)

//载荷发生器API
type Generator interface {
	Start() bool      //启动
	Stop() bool       //停止
	Status() uint32   //获取发生器状态
	CallCount() int64 //获取调用计数，每次启动都会重置计数
}

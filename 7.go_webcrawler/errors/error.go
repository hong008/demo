package errors

import (
	"bytes"
	"fmt"
	"strings"
)

type ErrorType string

const (
	ERROR_TYPE_DOWNLOADER ErrorType = "downloader error"
	ERROR_TYPE_ANALYZER   ErrorType = "analyzer error"
	ERROR_TYPE_PIPELINE   ErrorType = "pipeline error"
	ERROR_TYPE_SCHEDULER  ErrorType = "scheduler error"
)

type CrawlerError interface {
	Type() ErrorType //错误类型
	Error() string   //错误描述
}

type myCrawlerError struct {
	errType    ErrorType
	errMsg     string
	fullErrMsg string
}

func NewCrawlerError(errType ErrorType, errMsg string) CrawlerError {
	return &myCrawlerError{
		errType: errType,
		errMsg:  strings.TrimSpace(errMsg),
	}
}

func NewCrawlerErrorBy(errType ErrorType, err error) CrawlerError {
	return NewCrawlerError(errType, err.Error())
}

func (ce *myCrawlerError) Type() ErrorType {
	return ce.errType
}

func (ce *myCrawlerError) Error() string {
	if ce.fullErrMsg == "" {
		ce.getFullErrMsg()
	}
	return ce.fullErrMsg
}

func (ce *myCrawlerError) getFullErrMsg() {
	var buffer bytes.Buffer
	buffer.WriteString("crawler error: ")
	if ce.errType != "" {
		buffer.WriteString(string(ce.errType))
		buffer.WriteString(": ")
	}
	buffer.WriteString(ce.errMsg)
	ce.fullErrMsg = fmt.Sprintf("%s", buffer.String())
}

type IllegalParameterError struct {
	msg string
}

func NewIllegalParameterError(errMsg string) IllegalParameterError {
	return IllegalParameterError{
		msg:fmt.Sprintf("illegal parameter: %s", strings.TrimSpace(errMsg)),
	}
}

func (ipe IllegalParameterError) Error() string {
	return ipe.msg
}

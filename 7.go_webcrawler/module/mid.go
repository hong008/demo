package module

import (
	"demo/7.go_webcrawler/errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type MID string

var DefaultSNGen = NewSNGenertor(1, 0)

var midTemplate = "%s%s|%s"

//根据参数生成组件ID
func GenMID(mtype Type, sn uint64, maddr net.Addr) (MID, error) {
	if !LegalType(mtype) {
		errMsg := fmt.Sprintf("illegal module type: %s", mtype)
		return "", errors.NewIllegalParameterError(errMsg)
	}
	letter := legalTypeLetterMap[mtype]
	var midStr string
	if maddr == nil {
		midStr = fmt.Sprintf(midTemplate, letter, sn, "")
		midStr = midStr[:len(midStr)-1]
	} else {
		midStr = fmt.Sprintf(midTemplate, letter, sn, maddr.String())
	}
	return MID(midStr), nil
}

//判断给定的组件ID是否合法
func LegalMID(mid MID) bool {
	if _, err := SplitMID(mid); err == nil {
		return true
	}
	return false
}

//判断序列号是否合法
func legalSN(snStr string) bool {
	_, err := strconv.ParseUint(snStr, 10, 64)
	if err != nil {
		return false
	}
	return true
}

/*分解组件ID
第二个结果表示是否分解成功
如果分解成功，则第一个结果返回长度为3的string切片，内容分别为：
组件类型的字母，序列号，组件网络地址
*/
func SplitMID(mid MID) ([]string, error) {
	var ok bool
	var letter, snStr, addr string
	midStr := string(mid)

	if len(midStr) <= 1 {
		return nil, errors.NewIllegalParameterError("insufficient MID")
	}

	letter = midStr[:1]
	if _, ok = legalLetterTypeMap[letter]; !ok {
		return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module type letter: %s", letter))
	}

	snAndAddr := midStr[1:]
	index := strings.LastIndex(snAndAddr, "|")
	if index < 0 {
		snStr = snAndAddr
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module SN: %s", snStr))
		}
	} else {
		snStr = snAndAddr[:index]
		if !legalSN(snStr) {
			return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module SN: %s", snStr))
		}
		addr = snAndAddr[index+1:]
		index = strings.LastIndex(addr, ":")
		if index <= 0 {
			return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module addr: %s", addr))
		}
		ipStr := addr[:index]
		if ip := net.ParseIP(ipStr); ip == nil {
			return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module ip: %s", ip))
		}
		portStr := addr[index+1:]
		if _, err := strconv.ParseUint(portStr, 10, 64); err != nil {
			return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal module port: %s", portStr))
		}
	}
	return []string{letter, snStr, addr}, nil
}

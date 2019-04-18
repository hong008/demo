package module

import (
	"demo/7.go_webcrawler/errors"
	"fmt"
	"net"
	"strconv"
)

//组件网络地址
type mAddr struct {
	network string //网络协议
	address string //网络地址
}

func (maddr *mAddr) Network() string {
	return maddr.network
}

func (maddr *mAddr) String() string {
	return maddr.address
}

func NewAddr(network string, ip string, port uint64) (net.Addr, error) {
	if network != "http" && network != "https" {
		errMsg := fmt.Sprintf("illegal network for module address: %s", network)
		return nil, errors.NewIllegalParameterError(errMsg)
	}

	if parsedIp := net.ParseIP(ip); parsedIp == nil {
		return nil, errors.NewIllegalParameterError(fmt.Sprintf("illegal IP for module address: %s", ip))
	}

	return &mAddr{
		network: network,
		address: ip + ":" + strconv.Itoa(int(port)),
	}, nil
}

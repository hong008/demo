package __go_buffer

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type Buffer interface {
	Cap() uint32                         //获取缓冲器容量
	Len() uint32                         //获取缓冲器中的数据数量
	Put(datum interface{}) (bool, error) //向缓冲器放入数据
	Get() (interface{}, error)           //从缓冲器中获取数据
	Close() bool                         //关闭缓冲器
	Closed() bool                        //缓冲器是否已关闭
}

type myBuffer struct {
	ch          chan interface{} //数据通道
	closed      uint32           //缓冲器状态量：0-未关闭  1-已关闭
	closingLock sync.RWMutex     //为了消除因关闭缓冲器而产生的竞态条件的读写锁
}

func NewBuffer(size uint32) (Buffer, error) {
	if size == 0 {
		return nil, errors.New(fmt.Sprintf("illegal size for buffer: %d", size))
	}
	return &myBuffer{
		ch: make(chan interface{}, size),
	}, nil
}

func (buf *myBuffer) Cap() uint32 {
	return uint32(cap(buf.ch))
}

func (buf *myBuffer) Len() uint32 {
	return uint32(len(buf.ch))
}

func (buf *myBuffer) Put(datum interface{}) (ok bool, err error) {
	buf.closingLock.RLock()
	defer buf.closingLock.RUnlock()
	if buf.Closed() {
		return false, errors.New("closed buffer")
	}
	select {
	case buf.ch <- datum:
		ok = true
	default:
		ok = false

	}
	return
}

func (buf *myBuffer) Get() (interface{}, error) {
	select {
	case datum, ok := <-buf.ch:
		if !ok {
			return nil, errors.New("closed buffer")
		}
		return datum, nil
	default:
		return nil, nil
	}
}

func (buf *myBuffer) Close() bool {
	if atomic.CompareAndSwapUint32(&buf.closed, 0, 1) {
		buf.closingLock.Lock()
		close(buf.ch)
		buf.closingLock.Unlock()
		return true
	}
	return false
}

func (buf *myBuffer) Closed() bool {
	if atomic.LoadUint32(&buf.closed) == 0 {
		return false
	}
	return true
}

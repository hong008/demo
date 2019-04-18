package __go_buffer

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type Pool interface {
	BufferCap() uint32                   //获取缓冲器的统一容量
	MaxBufferNumber() uint32             //获取pool中缓冲器最大数量
	BufferNumber() uint32                //获取pool中缓冲器的数量
	Total() uint64                       //获取pool中数据的总数
	Put(datum interface{}) error         //向池中放入数据，若pool已关闭则直接返回错误
	Get() (datum interface{}, err error) //从pool中获取数据
	Close() bool
	Closed() bool
}

func NewPool(bufferCap uint32, maxBufferNumber uint32) (Pool, error) {
	if bufferCap == 0 {
		return nil, errors.New(fmt.Sprintf("illegal buffer cap for buffer pool: %d", bufferCap))
	}
	if maxBufferNumber == 0 {
		return nil, errors.New(fmt.Sprintf("illegal max buffer number for buffer pool: %d", maxBufferNumber))
	}
	bufCh := make(chan Buffer, maxBufferNumber)
	buf, _ := NewBuffer(bufferCap)
	bufCh <- buf
	return &myPool{
		bufferCap:       bufferCap,
		maxBufferNumber: maxBufferNumber,
		bufferNumber:    1,
		bufCh:           bufCh,
	}, nil
}

type myPool struct {
	bufferCap       uint32      //缓冲器统一容量
	maxBufferNumber uint32      //缓冲器的最大数量
	bufferNumber    uint32      //缓冲器实际数量
	total           uint64      //数据的总和
	bufCh           chan Buffer //存放缓冲器的通道
	closed          uint32      //pool的状态: 0-未关闭 1-已关闭
	rwlock          sync.RWMutex
}

func (pool *myPool) BufferCap() uint32 {
	return pool.bufferCap
}

func (pool *myPool) MaxBufferNumber() uint32 {
	return pool.maxBufferNumber
}

func (pool *myPool) BufferNumber() uint32 {
	return atomic.LoadUint32(&pool.bufferNumber)
}

func (pool *myPool) Total() uint64 {
	return atomic.LoadUint64(&pool.total)
}

func (pool *myPool) Put(datum interface{}) (err error) {
	if pool.Closed() {
		return errors.New("closed buffer pool")
	}
	var count uint32
	maxCount := pool.BufferNumber() * 5
	var ok bool
	for buf := range pool.bufCh {
		ok, err = pool.putData(buf, datum, &count, maxCount)
		if ok || err != nil {
			break
		}
	}
	return
}

func (pool *myPool) putData(buf Buffer, datum interface{}, count *uint32, maxCount uint32) (ok bool, err error) {
	if pool.Closed() {
		return false, errors.New("closed buffer pool")
	}

	defer func() {
		pool.rwlock.Lock()
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			err = errors.New("closed buffer pool")
		} else {
			pool.bufCh <- buf
		}
		pool.rwlock.RUnlock()
	}()

	ok, err = buf.Put(datum)
	if ok {
		atomic.AddUint64(&pool.total, 1)
	}
	if err != nil {
		return
	}
	(*count)++
	if *count >= maxCount && pool.BufferNumber() < pool.MaxBufferNumber() {
		pool.rwlock.Lock()
		if pool.BufferNumber() < pool.MaxBufferNumber() {
			if pool.Closed() {
				pool.rwlock.Unlock()
				return
			}
			newBuf, _ := NewBuffer(pool.bufferCap)
			newBuf.Put(datum)
			pool.bufCh <- newBuf
			atomic.AddUint32(&pool.bufferNumber, 1)
			atomic.AddUint64(&pool.total, 1)
			ok = true
		}
		pool.rwlock.Unlock()
		*count = 0
	}
	return
}

func (pool *myPool) Get() (datum interface{}, err error) {
	if pool.Closed() {
		return nil, errors.New("closed buffer pool")
	}
	var count uint32
	maxCount := pool.BufferNumber() * 10
	for buf := range pool.bufCh {
		datum, err = pool.getData(buf, &count, maxCount)
		if datum != nil || err != nil {
			break
		}
	}
	return
}

func (pool *myPool) getData(buf Buffer, count *uint32, maxCount uint32) (datum interface{}, err error) {
	if pool.Closed() {
		return nil, errors.New("closed buffer pool")
	}
	defer func() {
		if *count >= maxCount && buf.Len() == 0 && pool.BufferNumber() > 1 {
			buf.Close()
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			*count = 0
			return
		}
		pool.rwlock.RLock()
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			err = errors.New("closed buffer pool")
		} else {
			pool.bufCh <- buf
		}
		pool.rwlock.RUnlock()
	}()

	datum, err = buf.Get()
	if datum != nil {
		atomic.AddUint64(&pool.total, ^uint64(0))
		return
	}
	if err != nil {
		return
	}
	(*count)++
	return
}

func (pool *myPool) Close() bool {
	if !atomic.CompareAndSwapUint32(&pool.closed, 0, 1) {
		return false
	}
	pool.rwlock.Lock()
	defer pool.rwlock.Unlock()

	close(pool.bufCh)
	for buf := range pool.bufCh {
		buf.Close()
	}
	return true
}

func (pool *myPool) Closed() bool {
	if atomic.LoadUint32(&pool.closed) == 1 {
		return true
	}
	return false
}

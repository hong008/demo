package module

import (
	"math"
	"sync"
)

//组件序列号生成器
type SNGenertor interface {
	Start() uint64      //用于获取设置的最小的序列号
	Max() uint64        //获取设置的最大的序列号
	Next() uint64       //获取下一个序列号
	CycleCount() uint64 //序列号从最小到最大可以循环，此方法用于获取循环的次数
	Get() uint64        //获取序列号，并准备好下一个
}

//生成一个序列号生成器的实例
func NewSNGenertor(start, max uint64) SNGenertor {
	if max == 0 {
		max = math.MaxUint64
	}
	if max < start {
		return nil
	}
	return &mySNGenertor{
		start: start,
		max:   max,
		next:  start,
	}
}

//实现生成器接口的结构
type mySNGenertor struct {
	start      uint64       //序列号起始值
	max        uint64       //序列号最大值
	next       uint64       //下一个序列号
	cycleCount uint64       //循环次数
	lock       sync.RWMutex //锁
}

func (gen *mySNGenertor) Start() uint64 {
	return gen.start
}

func (gen *mySNGenertor) Max() uint64 {
	return gen.max
}

func (gen *mySNGenertor) Next() uint64 {
	gen.lock.RLock()
	defer gen.lock.RUnlock()
	return gen.next
}

func (gen *mySNGenertor) CycleCount() uint64 {
	gen.lock.RLock()
	defer gen.lock.RUnlock()
	return gen.cycleCount
}

func (gen *mySNGenertor) Get() uint64 {
	gen.lock.Lock()
	defer gen.lock.Unlock()

	id := gen.next
	if id == gen.max {
		gen.next = gen.start
		gen.cycleCount++
	} else {
		gen.cycleCount++
	}
	return id
}

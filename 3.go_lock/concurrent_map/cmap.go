package concurrent_map

import (
	"math"
	"sync/atomic"
)

type ConcurrentMap interface {
	Concurrency() int                                  //返回并发量
	Put(key string, element interface{}) (bool, error) //save一个键值对，值不能为nil，如果key已经存在，则直接替换旧值为新值
	Get(key string) interface{}                        //根据key获取一个value
	Delete(key string) bool                            //删除键值对
	Len() uint64                                       //返回键值对数量
}

type myConcurrentMap struct {
	concurrentcy int
	segments     []Segment
	total        uint64
}

func NewConcurrentMap(concurrency int, redistributor PairRedistributor) (ConcurrentMap, error) {
	if concurrency <= 0 {
		return nil, newIllegalParameterError("concurrency is too small")
	}
	if concurrency > MAX_CONCURRENCY {
		return nil, newIllegalParameterError("concurrency is to large")
	}
	cmap := &myConcurrentMap{}
	cmap.concurrentcy = concurrency
	cmap.segments = make([]Segment, concurrency)
	for i := 0; i < concurrency; i++ {
		cmap.segments[i] = newSegment(DEFAULT_BUCKET_NUMBER, redistributor)
	}
	return cmap, nil
}

func (cmap *myConcurrentMap) Concurrency() int {
	return cmap.concurrentcy
}

func (cmap *myConcurrentMap) Put(key string, element interface{}) (bool, error) {
	p, err := newPair(key, element)
	if err != nil {
		return false, err
	}
	s := cmap.findSegment(p.Hash())
	ok, err := s.Put(p)
	if ok {
		atomic.AddUint64(&cmap.total, 1)
	}
	return ok, err
}

func (cmap *myConcurrentMap) Get(key string) interface{} {
	keyHash := hash(key)
	s := cmap.findSegment(keyHash)
	pair := s.GetWithHash(key, keyHash)
	if pair == nil {
		return nil
	}
	return pair.Element()
}

func (cmap *myConcurrentMap) Delete(key string) bool {
	s := cmap.findSegment(hash(key))
	if s.Delete(key) {
		atomic.AddUint64(&cmap.total, ^uint64(0))
		return true
	}
	return false
}

func (cmap *myConcurrentMap) Len() uint64 {
	return atomic.LoadUint64(&cmap.total)
}

func (cmap *myConcurrentMap) findSegment(keyHash uint64) Segment {
	if cmap.concurrentcy == 1 {
		return cmap.segments[0]
	}
	var keyHash32 uint32
	if keyHash > math.MaxUint32 {
		keyHash32 = uint32(keyHash >> 32)
	} else {
		keyHash32 = uint32(keyHash)
	}
	return cmap.segments[int(keyHash32>>16)%(cmap.concurrentcy-1)]
}

package concurrent_map

type ConCurrentMap interface {
	ConCurrency() int                                  //返回并发量
	Put(key string, element interface{}) (bool, error) //save一个键值对，值不能为nil，如果key已经存在，则直接替换旧值为新值
	Get(key string) interface{}                        //根据key获取一个value
	Delete(key string) bool                            //删除键值对
	Len() uint64                                       //返回键值对数量
}

type myConCurrentMap struct {
	concurrentcy int
	segments []Segment
	total uint64
}

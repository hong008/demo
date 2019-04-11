package __go_lock

import "sync/atomic"

type ConCurrentArray interface {
	Set(index uint32, elem int) error
	Get(index uint32) (int, error)
	Len() uint32
}

type conCurrentArray struct {
	length uint32
	val    atomic.Value
}

func NewConCurrentArray(len uint32) ConCurrentArray {
	array := &conCurrentArray{}
	array.length = len
	array.val.Store(make([]int, array.length))
	return array
}

func (c *conCurrentArray) Set(index uint32, elem int) error {
	//todo check

	array := make([]int, c.length)
	array = c.val.Load().([]int)
	array[index] = elem
	c.val.Store(array)
	return nil
}

func (c *conCurrentArray) Get(index uint32) (val int, err error) {
	//todo check
	return c.val.Load().([]int)[index], nil
}

func (c *conCurrentArray) Len() uint32 {
	return c.length
}

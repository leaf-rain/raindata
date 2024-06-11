package slice

import (
	"sync"
)

type FixedSizeSlice struct {
	slice []interface{}
	size  int
	lock  sync.Mutex
}

func NewFixedSizeSlice(size int) *FixedSizeSlice {
	return &FixedSizeSlice{
		slice: make([]interface{}, 0, size),
		size:  size,
		lock:  sync.Mutex{},
	}
}

func (f *FixedSizeSlice) Push(val interface{}) interface{} {
	f.lock.Lock()
	defer f.lock.Unlock()
	var item interface{}
	if len(f.slice) == f.size {
		// 如果达到最大长度，移除第一个元素
		item = f.slice[0]
		f.slice = f.slice[1:]
	}
	f.slice = append(f.slice, val)
	return item
}

func (f *FixedSizeSlice) Pop() interface{} {
	f.lock.Lock()
	defer f.lock.Unlock()
	item := f.slice[0]
	f.slice = f.slice[1:]
	return item
}

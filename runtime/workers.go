package runtime

import (
	"sync"
	"sync/atomic"
)

type workersInfo struct{
	num int32
}

func (w *workersInfo) Add() {
	atomic.AddInt32(&w.num, 1)
}
func (w *workersInfo) Done() {
	atomic.AddInt32(&w.num, -1)
}

func (w workersInfo) Num() int32{
	return w.num
}

var (
	_workersInfoOnce sync.Once
	_workersInfo *workersInfo
)

func GetWorkersInfo() *workersInfo{
	_workersInfoOnce.Do(func() {
		_workersInfo = &workersInfo{num: 0}
	})
	return _workersInfo
}
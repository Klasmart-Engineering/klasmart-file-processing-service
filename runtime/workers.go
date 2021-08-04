package runtime

import (
	"sync"
	"sync/atomic"
)

type WorkersInfoData struct {
	Workers         int32               `json:"workers"`
	ProcessingFiles map[string][]string `json:"processing_files"`
}

type workersInfo struct {
	sync.RWMutex
	num             int32
	processingFiles map[string][]string
}

func (w *workersInfo) Add(topic string, message string) {
	w.Lock()
	defer w.Unlock()
	atomic.AddInt32(&w.num, 1)
	w.processingFiles[topic] = append(w.processingFiles[topic], message)
}
func (w *workersInfo) Done(topic string, message string) {
	w.Lock()
	defer w.Unlock()
	atomic.AddInt32(&w.num, -1)
	w.removeProcessingFile(topic, message)
}

func (w *workersInfo) Num() int32 {
	w.RLock()
	defer w.RUnlock()
	return w.num
}

func (w *workersInfo) Info() *WorkersInfoData {
	w.RLock()
	defer w.RUnlock()
	return &WorkersInfoData{
		Workers:         w.num,
		ProcessingFiles: w.processingFiles,
	}
}

func (w *workersInfo) removeProcessingFile(topic string, message string) {
	files := w.processingFiles[topic]
	if files != nil {
		//only one
		if len(files) <= 1 {
			w.processingFiles[topic] = nil
			return
		}

		index := -1
		for i := range files {
			if files[i] == message {
				index = i
				break
			}
		}
		if index >= 0 {
			newFiles := w.processingFiles[topic]
			newFiles = w.removeSlice(newFiles, index)
			w.processingFiles[topic] = newFiles
		}

	}
}
func (w *workersInfo) removeSlice(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
}

var (
	_workersInfoOnce sync.Once
	_workersInfo     *workersInfo
)

func GetWorkersInfo() *workersInfo {
	_workersInfoOnce.Do(func() {
		_workersInfo = &workersInfo{
			num:             0,
			processingFiles: make(map[string][]string),
		}
	})
	return _workersInfo
}

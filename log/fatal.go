package fatal

import (
	"context"
	"fmt"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"os"
	"sync"
	"time"
)

var (
	_failedFile *os.File
	_failedFileOnce sync.Once
)

func failedFile() *os.File {
	_failedFileOnce.Do(func() {
		var err error
		_failedFile, err = os.OpenFile(config.Get().Log.FailedFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	})
	return _failedFile
}

func Write(ctx context.Context, text string, values ...interface{}) {
	errMsg := fmt.Sprintf(text, values...)
	msg := fmt.Sprintf("[FATAL] %v %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
	failedFile().WriteString(msg)
}

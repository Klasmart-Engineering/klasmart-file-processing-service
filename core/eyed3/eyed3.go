package eyed3

import (
	"bytes"
	"context"
	"os/exec"
	"sync"

	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
)

const (
	defaultEyeD3Path = "./eyed3_py.py"
)

type EyeD3Tool struct {
	eyeD3Path string
}

func (e *EyeD3Tool) RemoveMP3MetaData(ctx context.Context, fileName string) error {
	args := []string{
		e.eyeD3Path,
		fileName,
	}
	cmd := exec.Command("python", args...)
	out := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = out
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		log.Error(ctx, "Run eyed3 failed",
			log.Err(err),
			log.String("stderr", stderr.String()),
			log.String("stdout", out.String()),
			log.String("file", fileName),
			log.String("python file", e.eyeD3Path))
	}
	return err
}
func (e *EyeD3Tool) SetEyeD3Path(path string) {
	if path == "" {
		path = defaultEyeD3Path
	}
	e.eyeD3Path = path
}

var (
	_eyeD3ToolOnce sync.Once
	_eyeD3Tool     *EyeD3Tool
)

func GetEyeD3Tool() *EyeD3Tool {
	_eyeD3ToolOnce.Do(func() {
		_eyeD3Tool = &EyeD3Tool{
			eyeD3Path: defaultEyeD3Path,
		}
		if config.Get().Core.EyeD3Path != "" {
			_eyeD3Tool.eyeD3Path = config.Get().Core.EyeD3Path
		}
	})
	return _eyeD3Tool
}

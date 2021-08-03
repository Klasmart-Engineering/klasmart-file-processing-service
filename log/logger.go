package log

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	_logger     *logrus.Logger
	_errFile    *os.File

	_loggerOnce     sync.Once
	_errFileOnce    sync.Once

	_failedFile *os.File
	_failedFileOnce sync.Once
)

func logName(prefix string) string {
	today := time.Now().Format("2006_01_02")
	return prefix + "_" + today + ".log"
}

func failedFile() *os.File {
	_failedFileOnce.Do(func() {
		var err error
		_failedFile, err = os.OpenFile("failed.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	})
	return _failedFile
}
func LoggerOut() io.Writer{
	return logger().Out
}
func ErrOut() io.Writer{
	return errFile()
}

func logger() *logrus.Logger {
	_loggerOnce.Do(func() {
		_logger = logrus.New()

		file, err := os.OpenFile(logName("log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		_logger.Out = file
		_logger.SetFormatter(&logrus.JSONFormatter{})
		//TODO: Maybe need to config
		_logger.SetLevel(logrus.DebugLevel)
	})
	return _logger
}

func errFile() *os.File {
	//open err file
	_errFileOnce.Do(func() {
		var err error
		_errFile, err = os.OpenFile(logName("err"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	})

	return _errFile
}

func Debug(ctx context.Context, text string, values ...interface{}) {
	logger().WithField("context", ctx).Debug(fmt.Sprintf(text, values...))
	printLog(ctx, "DEBUG", text, values...)
}

func Info(ctx context.Context, text string, values ...interface{}) {
	logger().WithField("context", ctx).Info(fmt.Sprintf(text, values...))
	printLog(ctx, "INFO", text, values...)
}

func Warn(ctx context.Context, text string, values ...interface{}) {
	logger().WithField("context", ctx).Warn(fmt.Sprintf(text, values...))
	printLog(ctx, "WARN", text, values...)
}

func Error(ctx context.Context, text string, values ...interface{}) {
	errMsg := fmt.Sprintf(text, values...)
	logger().WithField("context", ctx).Error(errMsg)
	msg := fmt.Sprintf("[ERROR] %v %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
	errFile().WriteString(msg)

	printLog(ctx, "ERROR", text, values...)
}

func Failed(ctx context.Context, text string, values ...interface{}) {
	errMsg := fmt.Sprintf(text, values...)
	msg := fmt.Sprintf("[FATAL] %v %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
	failedFile().WriteString(msg)

	printLog(ctx, "FATAL", text, values...)
}

func printLog(ctx context.Context, level, text string, values ...interface{}) {
	if !config.Get().Log.StdOut {
		return
	}
	msg := fmt.Sprintf(text, values...)
	out := fmt.Sprintf("[%v] %v %s ", level, time.Now().Format("2006-01-02 15:04:05"), msg)
	if !strings.HasSuffix(out, "\n") {
		out = out + "\n"
	}
	fmt.Print(out)
}
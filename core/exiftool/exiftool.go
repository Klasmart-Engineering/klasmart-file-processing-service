package exiftool

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"io"
	"os/exec"
	"runtime"
	"sync"
)

const (
	defaultExifToolPath = "exiftool"
)

var (
	initArgs   = []string{"-stay_open", "true", "-@", "-"}
	closeArgs  = []string{"-stay_open", "false"}
	commonArgs = []string{"-j"}
)

type ExifTool struct {
	exifToolPath string
	cmd          *exec.Cmd
	stdout       io.ReadCloser
	stdin        io.WriteCloser
	scanner      *bufio.Scanner
}

func (e *ExifTool) Start() error {
	args := make([]string, 0, len(initArgs))
	args = append(args, initArgs...)
	if len(commonArgs) > 0 {
		args = append(args, "-common_args")
		args = append(args, commonArgs...)
	}

	e.cmd = exec.Command(e.exifToolPath, args...)
	r, w := io.Pipe()
	e.stdout = r

	e.cmd.Stdout = w
	e.cmd.Stderr = w

	var err error
	e.stdin, err = e.cmd.StdinPipe()
	if err != nil {
		return err
	}

	e.scanner = bufio.NewScanner(r)
	e.scanner.Split(splitReadyToken)

	return e.cmd.Start()
}

func (e *ExifTool) Stop() {
	//close exiftool
	e.addArgs(closeArgs)
	e.execute()
	//close io
	e.stdout.Close()
	e.stdin.Close()
}
func (e *ExifTool) RemoveMetadata(ctx context.Context, file, outFile string, tags []string) error {
	//run tag params
	args := make([]string, 0, len(tags))
	args = append(args, e.tagsToArgs(tags)...)
	//add input & output files
	args = append(args, file)
	if outFile != "" {
		args = append(args, "-o")
		args = append(args, outFile)
	}

	//run command
	out, err := e.runCommand(ctx, args)

	//read result
	if err != nil {
		log.Error(ctx, "execute metadata failed",
			log.Err(err),
			log.String("input", file),
			log.String("output", outFile),
			log.String("stdout", out),
			log.Strings("tags", tags))
		return e.scanner.Err()
	}
	log.Info(ctx, "execute metadata successfully",
		log.String("input", file),
		log.String("output", outFile),
		log.String("out", out),
		log.Strings("tags", tags))
	return nil
}

func (e *ExifTool) runCommand(ctx context.Context, args []string) (string, error) {
	//add commands
	e.addArgs(args)

	//run command
	e.execute()

	if !e.scanner.Scan() {
		log.Error(ctx, "nothing on stdMergedOut",
			log.Err(e.scanner.Err()),
			log.Strings("args", args))
		return "", e.scanner.Err()
	}

	//read result
	if e.scanner.Err() != nil {
		log.Error(ctx, "execute remove metadata failed",
			log.Err(e.scanner.Err()),
			log.Strings("args", args))
		return "", e.scanner.Err()
	}
	log.Info(ctx, "execute remove metadata successfully",
		log.String("stdout", e.scanner.Text()),
		log.Strings("tags", args))
	return e.scanner.Text(), nil
}

func (e *ExifTool) SetExifToolPath(path string) {
	if path == "" {
		path = defaultExifToolPath
	}
	e.exifToolPath = path
}

func (e *ExifTool) execute() {
	fmt.Fprintf(e.stdin, "-execute\n")
}
func (e *ExifTool) addArgs(args []string) {
	for i := range args {
		fmt.Fprintf(e.stdin, args[i]+"\n")
	}
}

func (e ExifTool) tagToArg(tag string) string {
	return "-" + tag + "="
}
func (e ExifTool) tagsToArgs(tags []string) []string {
	args := make([]string, len(tags))
	for i := range tags {
		args[i] = e.tagToArg(tags[i])
	}
	return args
}

func splitReadyToken(data []byte, atEOF bool) (int, []byte, error) {
	readyToken := []byte("{ready}\n")
	if runtime.GOOS == "windows" {
		readyToken = []byte("{ready}\r\n")
	}

	var readyTokenLen = len(readyToken)

	idx := bytes.Index(data, readyToken)
	if idx == -1 {
		if atEOF && len(data) > 0 {
			return 0, data, fmt.Errorf("no final token found")
		}
		return 0, nil, nil
	}
	return idx + readyTokenLen, data[:idx], nil
}

var (
	_exifToolOnce sync.Once
	_exifTool     *ExifTool
)

func GetExifTool() *ExifTool {
	_exifToolOnce.Do(func() {
		_exifTool = &ExifTool{
			exifToolPath: defaultExifToolPath,
		}
	})
	return _exifTool
}

package entity

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/log"
	"os"
	"strings"
)

const(
	MQPrefix = "kfps:"
)

type FileInfo struct {
	Classify string
	Extension string
	Name string
	Path string
}
func (f FileInfo) Topic() string {
	return ""
}

func ParseFileInfo(topic, message string) *FileInfo {
	if !strings.HasPrefix(topic, MQPrefix) {
		return nil
	}
	classify := topic[len(MQPrefix):]

	namePairs := strings.Split(message, "/")
	extensionPairs := strings.Split(message, ".")
	name := namePairs[len(namePairs) - 1]
	extension := extensionPairs[len(extensionPairs) - 1]

	return &FileInfo{
		Classify:  classify,
		Extension: extension,
		Name:      name,
		Path:      message,
	}
}

type HandleFileParams struct {
	Classify string
	Extension string
	Name      string
	LocalFile *os.File
	LocalPath string

	DistFile *os.File
	DistPath string
}
func (h *HandleFileParams) CreateOutputFile(ctx context.Context) (*os.File, error){
	path := os.TempDir() + "/" + h.Name + "-handled"
	dst, err := os.Create(path)
	if err != nil {
		log.Error(ctx, "Can't create dist file, err: ", err)
		return nil, err
	}

	h.DistFile = dst
	h.DistPath = path
	return dst, nil
}
func (h *HandleFileParams) CleanOutputFile(ctx context.Context) {
	if h.DistFile != nil {
		h.DistFile.Close()
		err := os.Remove(h.DistPath)
		if err != nil {
			log.Error(ctx,"Can't remove output file, err: %v", err)
		}
	}
}
func (h *HandleFileParams) CleanLocalFile(ctx context.Context) {
	if h.LocalFile != nil {
		h.LocalFile.Close()
		err := os.Remove(h.LocalPath)
		if err != nil {
			log.Error(ctx,"Remove local file failed, err: %v", err)
		}
	}
}
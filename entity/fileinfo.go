package entity

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"os"
	"strings"
)

type FileInfo struct {
	Extension string
	Name      string
	Path      string
}

func ParseFileInfo(message string) *FileInfo {

	namePairs := strings.Split(message, "/")
	extensionPairs := strings.Split(message, ".")
	name := namePairs[len(namePairs)-1]
	extension := extensionPairs[len(extensionPairs)-1]

	return &FileInfo{
		Extension: extension,
		Name:      name,
		Path:      message,
	}
}

type HandleFileParams struct {
	Classify  string
	Extension string
	Name      string
	LocalFile *os.File
	LocalPath string

	DistFile *os.File
	DistPath string
}

func (h *HandleFileParams) OutputFilePath(ctx context.Context) string {
	path := os.TempDir() + string(os.PathSeparator) + h.Name + "-handled"
	h.DistPath = path
	return path
}
func (h *HandleFileParams) CreateOutputFile(ctx context.Context) (*os.File, error) {
	path := h.OutputFilePath(ctx)
	dst, err := os.Create(path)
	if err != nil {
		log.Error(ctx, "Can't create dist file",
			log.Err(err),
			log.Any("params", h))
		return nil, err
	}

	h.DistFile = dst
	return dst, nil
}
func (h *HandleFileParams) CleanOutputFile(ctx context.Context) {
	if h.DistFile != nil {
		h.DistFile.Close()
	}
	err := os.Remove(h.DistPath)
	if err != nil {
		log.Error(ctx, "Can't remove output file",
			log.Err(err),
			log.Any("params", h))
	}
}
func (h *HandleFileParams) CleanLocalFile(ctx context.Context) {
	if h.LocalFile != nil {
		h.LocalFile.Close()
		err := os.Remove(h.LocalPath)
		if err != nil {
			log.Error(ctx, "Remove local file failed",
				log.Err(err),
				log.Any("params", h))
		}
	}
}

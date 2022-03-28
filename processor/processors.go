package processor

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"strings"
	"sync"

	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
)

type IFileProcessor interface {
	HandleFile(ctx context.Context, f *entity.HandleFileParams) error
	SupportExtensions() []string
}

type ExifProcessor struct {
}

func (a ExifProcessor) HandleFile(ctx context.Context, f *entity.HandleFileParams) error {
	switch strings.ToLower(f.Extension) {
	case "jpg":
		return core.GetRemoveJPEGMetaDataHandler().Do(ctx, f)
	case "jpeg":
		return core.GetRemoveJPEGMetaDataHandler().Do(ctx, f)
	case "mp4":
		return core.GetRemoveMP4MetaDataHandler().Do(ctx, f)
	case "mov":
		return core.GetRemoveMOVMetaDataHandler().Do(ctx, f)
	case "mp3":
		return core.GetRemoveMP3MetaDataHandler().Do(ctx, f)
	case "png":
		return nil
	case "git":
		return nil
	case "bmp":
		return nil
	case "docx", "xlsx", "pptx":
		return core.GetRemoveOOXMLMetaDataHandler().Do(ctx, f)
	}
	return nil
}

func (a ExifProcessor) SupportExtensions() []string {
	return []string{
		"jpg", "jpeg", "mp4", "mp3", "mov", "docx", "xlsx", "pptx",
	}
}

var (
	_attachmentProcessor     *ExifProcessor
	_attachmentProcessorOnce sync.Once
)

func GetProcessor(p string) IFileProcessor {
	switch p {
	case "exif":
		_attachmentProcessorOnce.Do(func() {
			_attachmentProcessor = new(ExifProcessor)
		})
		return _attachmentProcessor
	default:
		log.Info(nil, "No processor found with key ",
			log.String("file", p))
	}
	return nil
}

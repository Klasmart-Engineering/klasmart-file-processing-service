package processor

import (
	"context"
	"strings"
	"sync"

	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
)

type IFileProcessor interface {
	HandleFile(ctx context.Context, f *entity.HandleFileParams) error
	SupportExtensions() []string
}

type AttachmentProcessor struct {
}

func (a AttachmentProcessor) HandleFile(ctx context.Context, f *entity.HandleFileParams) error {
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
	}
	return nil
}

func (a AttachmentProcessor) SupportExtensions() []string {
	return []string{
		"jpg", "jpeg", "mp4", "mp3", "mov",
	}
}

var (
	_attachmentProcessor     *AttachmentProcessor
	_attachmentProcessorOnce sync.Once
)

func GetAttachmentProcessor() IFileProcessor {
	_attachmentProcessorOnce.Do(func() {
		_attachmentProcessor = new(AttachmentProcessor)
	})
	return _attachmentProcessor
}

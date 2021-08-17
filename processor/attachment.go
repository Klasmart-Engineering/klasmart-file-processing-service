package processor

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
	"sync"
)

type IFileProcessor interface {
	HandleFile(ctx context.Context, f *entity.HandleFileParams) error
	SupportExtensions() []string
}

type AttachmentProcessor struct {
}

func (a AttachmentProcessor) HandleFile(ctx context.Context, f *entity.HandleFileParams) error {
	switch f.Extension {
	case "jpg":
		return core.GetRemoveJPEGMetaDataHandler().Do(ctx, f)
	case "jpeg":
		return core.GetRemoveJPEGMetaDataHandler().Do(ctx, f)
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
		"jpg", "png",
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

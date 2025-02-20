package core

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core/exiftool"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
	"sync"
)

type IFileHandler interface {
	Do(ctx context.Context, f *entity.HandleFileParams) error
}

type RemoveJPEGMetaDataHandler struct {
}

func (ih *RemoveJPEGMetaDataHandler) Do(ctx context.Context, f *entity.HandleFileParams) error {
	distPath := f.OutputFilePath(ctx)
	//_, err := f.CreateOutputFile(ctx)
	err := exiftool.GetExifTool().RemoveMetadata(ctx, f.LocalPath, distPath, exiftool.JpegTags)
	if err != nil {
		log.Error(ctx, "RemoveMetadata failed",
			log.Err(err),
			log.Any("params", f),
			log.Strings("tags", exiftool.JpegTags))
		return err
	}

	return nil
}

var (
	_removeJPEGMetaDataHandler     IFileHandler
	_removeJPEGMetaDataHandlerOnce sync.Once
)

func GetRemoveJPEGMetaDataHandler() IFileHandler {
	_removeJPEGMetaDataHandlerOnce.Do(func() {
		_removeJPEGMetaDataHandler = new(RemoveJPEGMetaDataHandler)
	})
	return _removeJPEGMetaDataHandler
}

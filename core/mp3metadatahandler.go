package core

import (
	"context"
	"sync"

	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core/exiftool"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core/eyed3"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
)

type RemoveMP3MetaDataHandler struct {
}

func (ih *RemoveMP3MetaDataHandler) Do(ctx context.Context, f *entity.HandleFileParams) error {
	err := eyed3.GetEyeD3Tool().RemoveMP3MetaData(ctx, f.LocalPath)
	if err != nil {
		log.Error(ctx, "RemoveMetadata failed",
			log.Err(err),
			log.Any("params", f),
			log.Strings("tags", exiftool.JpegTags))
		return err
	}
	f.DistPath = f.LocalPath

	return nil
}

var (
	_removeMP3MetaDataHandler     IFileHandler
	_removeMP3MetaDataHandlerOnce sync.Once
)

func GetRemoveMP3MetaDataHandler() IFileHandler {
	_removeMP3MetaDataHandlerOnce.Do(func() {
		_removeMP3MetaDataHandler = new(RemoveMP3MetaDataHandler)
	})
	return _removeMP3MetaDataHandler
}

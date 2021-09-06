package core

import (
	"context"
	"sync"

	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core/exiftool"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
)

type RemoveMOVMetaDataHandler struct {
}

func (ih *RemoveMOVMetaDataHandler) Do(ctx context.Context, f *entity.HandleFileParams) error {
	distPath := f.OutputFilePath(ctx)
	//_, err := f.CreateOutputFile(ctx)
	err := exiftool.GetExifTool().RemoveMetadata(ctx, f.LocalPath, distPath, exiftool.MovTags)
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
	_removeMOVMetaDataHandler     IFileHandler
	_removeMOVMetaDataHandlerOnce sync.Once
)

func GetRemoveMOVMetaDataHandler() IFileHandler {
	_removeMOVMetaDataHandlerOnce.Do(func() {
		_removeMOVMetaDataHandler = new(RemoveMOVMetaDataHandler)
	})
	return _removeMOVMetaDataHandler
}

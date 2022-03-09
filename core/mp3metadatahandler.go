package core

import (
	"context"
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

func GetRemoveMP3MetaDataHandler() IFileHandler {
	return new(RemoveMP3MetaDataHandler)
}

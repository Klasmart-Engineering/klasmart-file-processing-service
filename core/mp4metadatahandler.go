package core

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core/exiftool"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
)

type RemoveMP4MetaDataHandler struct {
}

func (ih *RemoveMP4MetaDataHandler) Do(ctx context.Context, f *entity.HandleFileParams) error {
	distPath := f.OutputFilePath(ctx)
	//_, err := f.CreateOutputFile(ctx)
	err := exiftool.GetExifTool().RemoveMetadata(ctx, f.LocalPath, distPath, exiftool.Mp4Tags)
	if err != nil {
		log.Error(ctx, "RemoveMetadata failed",
			log.Err(err),
			log.Any("params", f),
			log.Strings("tags", exiftool.JpegTags))
		return err
	}

	return nil
}

func GetRemoveMP4MetaDataHandler() IFileHandler {
	return new(RemoveMP4MetaDataHandler)
}

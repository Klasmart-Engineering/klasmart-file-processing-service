package core

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core/mp3"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
)

type RemoveMP3MetaDataHandler struct {
}

func (ih *RemoveMP3MetaDataHandler) Do(ctx context.Context, f *entity.HandleFileParams) error {
	distPath := f.OutputFilePath(ctx)
	err := mp3.RemoveMetadata(ctx, f.LocalPath, distPath)
	if err != nil {
		log.Error(ctx, "RemoveMetadata failed",
			log.Err(err),
			log.Any("params", f))
		return err
	}
	return nil
}

func GetRemoveMP3MetaDataHandler() IFileHandler {
	return new(RemoveMP3MetaDataHandler)
}

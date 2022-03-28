package core

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
	"image/png"
)

type RemovePNGMetaDataHandler struct {
}

func (ih *RemovePNGMetaDataHandler) Do(ctx context.Context, f *entity.HandleFileParams) error {
	dst, err := f.CreateOutputFile(ctx)
	if err != nil {
		log.Error(ctx, "Can't create output file",
			log.Err(err),
			log.Any("params", f))
		return err
	}

	img, err := png.Decode(f.LocalFile) //Decode file
	if err != nil {
		log.Error(ctx, "Can't decode png file",
			log.Err(err),
			log.Any("params", f))
		return err
	}
	err = png.Encode(dst, img) //Encode file
	if err != nil {
		log.Error(ctx, "Can't encode png file",
			log.Err(err),
			log.Any("params", f))
		return err
	}

	return nil
}

func GetRemovePNGMetaDataHandler() IFileHandler {
	return new(RemovePNGMetaDataHandler)
}

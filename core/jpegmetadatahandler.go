package core

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/log"
	"image/jpeg"
	"sync"
)
type IFileHandler interface{
	Do(ctx context.Context, f *entity.HandleFileParams) error
}

type RemoveJPEGMetaDataHandler struct {
}

func (ih *RemoveJPEGMetaDataHandler) Do(ctx context.Context, f *entity.HandleFileParams) error {
	dst, err := f.CreateOutputFile(ctx)
	if err != nil {
		log.Error(ctx, "Can't create output file, err: ", err)
		return err
	}

	img, err := jpeg.Decode(f.LocalFile) //Decode file
	if err != nil {
		log.Error(ctx, "Can't decode jpeg file, err: ", err)
		return err
	}
	err = jpeg.Encode(dst, img, &jpeg.Options{Quality: 100}) //Encode file
	if err != nil {
		log.Error(ctx, "Can't encode jpeg file, err: ", err)
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

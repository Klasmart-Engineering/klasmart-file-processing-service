package service

import (
	"context"
	"errors"
	"fmt"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
)

type FileProcessingService struct {
	handler              map[string]func(ctx context.Context, f *entity.HandleFileParams) error
	supportExtensionsMap map[string][]string
}

func (fp *FileProcessingService) Handle(ctx context.Context, file string) error {

	err := fp.initProcessors()
	if err != nil {
		return err
	}

	for key, handler := range fp.handler {
		err := fp.handleMessage(ctx, file, key, handler)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fp *FileProcessingService) SupportExtensions() map[string][]string {
	return fp.supportExtensionsMap
}

func (fp *FileProcessingService) handleMessage(ctx context.Context,
	file, key string,
	handler func(ctx context.Context, f *entity.HandleFileParams) error) error {
	//parse file info

	fileInfo := entity.ParseFileInfo(file)
	if fileInfo == nil {
		log.Info(ctx, "parseFileInfo failed",
			log.String("file", file))
		return errors.New("failed to parse info file: " + file)
	}
	fmt.Println(fp.supportExtensionsMap)
	log.Info(ctx, "Check contains",
		log.String("fileInfo.Extension", fileInfo.Extension),
		log.Strings("fp.supportExtensionsMap[key]", fp.supportExtensionsMap[key]))
	//ignore unsupported extension
	supportExtension := fp.containsString(fileInfo.Extension, fp.supportExtensionsMap[key])
	if !supportExtension {
		log.Info(ctx, "Unsupported extension",
			log.String("fileInfo.Extension", fileInfo.Extension),
			log.Strings("fp.supportExtensionsMap[key]", fp.supportExtensionsMap[key]))
		return nil
	}

	log.Info(ctx, "downloading File",
		log.Any("fileInfo", fileInfo))
	//download file
	fileParams, err := fp.downloadFile(ctx, fileInfo)
	if err != nil {
		log.Error(ctx, "downloadFile failed",
			log.Err(err),
			log.Any("fileInfo", fileInfo))
		return err
	}
	log.Debug(ctx, "downloading success",
		log.Any("fileParams", fileParams))
	defer fileParams.CleanLocalFile(ctx)
	defer fileParams.CleanOutputFile(ctx)

	//handle file
	err = handler(ctx, fileParams)
	if err != nil {
		log.Error(ctx, "Handle file failed",
			log.Err(err),
			log.Any("fileParams", fileParams))
		return err
	}

	//upload file
	err = fp.uploadHandledFile(ctx, fileInfo, fileParams.DistPath)
	if err != nil {
		log.Error(ctx, "uploadHandledFile failed",
			log.Err(err),
			log.Any("fileInfo", fileInfo),
			log.Any("source", fileInfo),
			log.Any("fileParams", fileParams))
		return err
	}

	log.Debug(ctx, "handle file success")
	return nil
}

func GetFileProcessingService() *FileProcessingService {
	return &FileProcessingService{
		supportExtensionsMap: make(map[string][]string),
		handler:              make(map[string]func(ctx context.Context, f *entity.HandleFileParams) error),
	}
}

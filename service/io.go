package service

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/storage"
	"io"
	"os"
)

func (fp *FileProcessingService) removeUnusedFiles(ctx context.Context, filePath []string) {
	for i := range filePath {
		if filePath[i] == "" {
			continue
		}
		err := os.RemoveAll(filePath[i])
		if err != nil {
			log.Error(ctx, "Remove file failed",
				log.Err(err),
				log.Strings("filePath", filePath))
		}
	}
}

func (fp *FileProcessingService) uploadHandledFile(ctx context.Context, fileInfo *entity.FileInfo, filePath string) error {
	if filePath == "" {
		return nil
	}
	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
	if err != nil {
		log.Error(ctx, "Can't open fileInfo",
			log.Err(err),
			log.String("filePath", filePath),
			log.Any("fileInfo", fileInfo))
		return err
	}
	defer f.Close()

	uploadPath := fileInfo.Path
	err = storage.DefaultStorage(ctx).UploadFile(ctx, uploadPath, f)
	if err != nil {
		log.Error(ctx, "Can't upload resource",
			log.Err(err),
			log.Any("f", f),
			log.String("uploadPath", uploadPath))
		return err
	}
	return nil
}

func (fp *FileProcessingService) backupFile(ctx context.Context, file *entity.FileInfo, f *os.File) error {
	uploadPath := "/backup/" + file.Path
	err := storage.DefaultStorage(ctx).UploadFile(ctx, uploadPath, f)
	if err != nil {
		log.Error(ctx, "Can't upload resource",
			log.Err(err),
			log.Any("f", f),
			log.Any("file", file))
		return err
	}
	f.Seek(0, io.SeekStart)
	return nil
}

func (fp *FileProcessingService) downloadFile(ctx context.Context, fileInfo *entity.FileInfo) (*entity.HandleFileParams, error) {
	reader, err := storage.DefaultStorage(ctx).DownloadFile(ctx, fileInfo.Path)
	if err != nil {
		log.Error(ctx, "Can't download resource",
			log.Err(err),
			log.Any("fileInfo", fileInfo))
		return nil, err
	}

	localFilePath := os.TempDir() + "/" + fileInfo.Name

	f, err := os.OpenFile(localFilePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Error(ctx, "Can't open fileInfo",
			log.Err(err),
			log.Any("fileInfo", fileInfo))
		return nil, err
	}
	_, err = io.Copy(f, reader)
	if err != nil {
		log.Error(ctx, "Save fileInfo failed",
			log.Err(err),
			log.Any("f", f))
		return nil, err
	}
	f.Close()

	//Open fileInfo with read only mode
	f, err = os.OpenFile(localFilePath, os.O_RDONLY, 0777)
	if err != nil {
		log.Error(ctx, "Can't open fileInfo",
			log.Err(err),
			log.String("localFilePath", localFilePath))
		return nil, err
	}

	return &entity.HandleFileParams{
		Extension: fileInfo.Extension,
		Name:      fileInfo.Name,
		Classify: fileInfo.Classify,
		LocalFile: f,
		LocalPath: localFilePath,
	}, nil
}

func (fp *FileProcessingService) containsString(extension string, supportedExtensions []string) bool{
	for i := range supportedExtensions {
		if supportedExtensions[i] == extension {
			return true
		}
	}
	return false
}

package storage

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"io"
	"mime/multipart"
	"sync"
)

var (
	_defaultStorageOnce sync.Once
	_defaultStorage     IStorage
)

type IStorage interface {
	OpenStorage(ctx context.Context) error
	CloseStorage(ctx context.Context)
	UploadFile(ctx context.Context, filePath string, fileStream multipart.File) error
	DownloadFile(ctx context.Context, filePath string) (io.Reader, error)
	ExistFile(ctx context.Context, filePath string) (int64, bool)

	ListAll() ([]string, error)
}

//根据环境变量创建存储对象
func createStorageByEnv(ctx context.Context) {
	conf := config.Get()

	switch conf.Storage.Driver {
	case "s3":
		_defaultStorage = newS3Storage(S3StorageConfig{
			Endpoint:   conf.Storage.EndPoint,
			Bucket:     conf.Storage.Bucket,
			BucketOut:  conf.Storage.BucketOut,
			Region:     conf.Storage.Region,
			AWSSession: conf.Storage.AWSSession,
			//SecretID:   conf.Storage.SecretID,
			//SecretKey:  conf.Storage.SecretKey,
			Accelerate: conf.Storage.Accelerate,
		})
		err := _defaultStorage.OpenStorage(ctx)
		if err != nil {
			log.Error(ctx, "open storage failed",
				log.Err(err),
				log.Any("config", conf.Storage))
			panic(err)
		}
	default:
		panic("Environment CLOUD_ENV is nil")
	}
}
func DefaultStorage(ctx context.Context) IStorage {
	_defaultStorageOnce.Do(func() {
		if _defaultStorage == nil {
			createStorageByEnv(ctx)
		}
	})
	return _defaultStorage
}

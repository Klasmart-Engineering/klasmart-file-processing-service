package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
)

type S3StorageConfig struct {
	Bucket     string
	BucketOut  string
	Region     string
	Accelerate bool
}

type S3Storage struct {
	session    *session.Session
	bucket     string
	bucketOut  string
	region     string
	accelerate bool
}

func (s *S3Storage) OpenStorage(ctx context.Context) error {
	cfg := &aws.Config{
		Region:          aws.String(s.region),
		S3UseAccelerate: aws.Bool(s.accelerate),
	}
	sess, err := session.NewSession(cfg)
	if err != nil {
		return err
	}

	s.session = sess
	return nil
}
func (s *S3Storage) CloseStorage(ctx context.Context) {

}

func (s *S3Storage) UploadFile(ctx context.Context, filePath string, fileStream multipart.File) error {
	uploader := s3manager.NewUploader(s.session)

	extension, err := s.fetchFileContentType(ctx, filePath)
	if err != nil {
		log.Error(ctx, "Fetch extension failed",
			log.Err(err),
			log.String("key", filePath))
		return err
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.bucketOut),
		Key:         aws.String(filePath),
		Body:        fileStream,
		ContentType: aws.String(extension),
	})
	if err != nil {
		log.Error(ctx, "File upload failed",
			log.Err(err),
			log.String("bucket", s.bucketOut),
			log.String("key", filePath))
		return err
	}
	return nil
}

func (s *S3Storage) DownloadFile(ctx context.Context, filePath string) (io.Reader, error) {
	downloader := s3manager.NewDownloader(s.session)
	data := make([]byte, 1024)
	writerAt := aws.NewWriteAtBuffer(data)

	_, err := downloader.Download(writerAt, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})

	if err != nil {
		log.Error(ctx, "download resource failed", log.Err(err))
		return nil, err
	}

	buffer := bytes.NewReader(writerAt.Bytes())
	return buffer, nil
}

func (s *S3Storage) ExistFile(ctx context.Context, filePath string) (int64, bool) {
	svc := s3.New(s.session)
	res, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})

	if err != nil {
		return -1, false
	}
	return *res.ContentLength, true
}

func (s *S3Storage) fetchFileContentType(ctx context.Context, key string) (string, error) {
	fileNamePairs := strings.Split(key, ".")
	if len(fileNamePairs) < 2 {
		return "", errors.New("no extension")
	}
	extension := fileNamePairs[len(fileNamePairs)-1]

	switch strings.ToLower(extension) {
	case "mp4":
		return "video/mpeg4", nil
	case "pdf":
		return "application/pdf", nil
	case "avi":
		return "video/avi", nil
	case "jpg":
		return "image/jpeg", nil
	case "jpeg":
		return "image/jpeg", nil
	case "jfif":
		return "image/jpeg", nil
	case "m4v":
		return "video/x-m4v", nil
	case "m4a":
		return "audio/m4a", nil
	case "webm":
		return "video/webm", nil
	case "gif":
		return "image/gif", nil
	case "png":
		return "image/png", nil
	case "mp3":
		return "audio/mp3", nil
	case "docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document", nil
	case "pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation", nil
	case "doc":
		return "application/msword", nil
	case "ppt":
		return "application/vnd.ms-powerpoint", nil
	case "xls":
		return "application/vnd.ms-excel", nil
	case "xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
	case "wav":
		return "audio/wav", nil
	case "mov":
		return "video/x-sgi-movie", nil
	case "pps":
		return "application/vnd.ms-powerpoint", nil
	case "ppsx":
		return "application/vnd.ms-powerpoint", nil
	}
	return "", errors.New("unknown file extension")
}

func newS3Storage(c S3StorageConfig) IStorage {
	return &S3Storage{
		bucket:     c.Bucket,
		bucketOut:  c.BucketOut,
		region:     c.Region,
		accelerate: c.Accelerate,
	}
}

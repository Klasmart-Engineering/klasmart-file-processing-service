package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type S3StorageConfig struct {
	Endpoint   string
	Bucket     string
	Region     string
	Accelerate bool
	SecretID string
	SecretKey string
}

type S3Storage struct {
	session    *session.Session
	bucket     string
	region     string
	endpoint   string
	accelerate bool

	secretID string
	secretKey string
}

type EndPointWithScheme struct {
	endpoint *string
	scheme   string
	isHttps  bool
}

func (s S3Storage) getEndpoint(ctx context.Context) (*EndPointWithScheme, error) {
	if s.endpoint == "" {
		return &EndPointWithScheme{
			endpoint: nil,
			scheme:   "https",
			isHttps:  true,
		}, nil
	}
	p, err := url.Parse(s.endpoint)
	if err != nil {
		return nil, err
	}
	ret := &EndPointWithScheme{
		endpoint: aws.String(s.endpoint),
		scheme:   p.Scheme,
		isHttps:  p.Scheme == "https",
	}

	return ret, nil
}

func (s *S3Storage) OpenStorage(ctx context.Context) error {
	//在~/.aws/credentials文件中保存secretId和secretKey
	endPointInfo, err := s.getEndpoint(ctx)
	if err != nil {
		return err
	}
	flag := !endPointInfo.isHttps

	sess, err := session.NewSession(&aws.Config{
		Endpoint:         endPointInfo.endpoint,
		Region:           aws.String(s.region),
		S3UseAccelerate:  aws.Bool(s.accelerate),
		DisableSSL:       aws.Bool(flag),
		S3ForcePathStyle: aws.Bool(flag),
		Credentials: 	credentials.NewStaticCredentials(s.secretID, s.secretKey, ""),
	})
	if err != nil {
		return err
	}

	s.session = sess
	return nil
}
func (s *S3Storage) CloseStorage(ctx context.Context) {

}

func getContentType(fileStream multipart.File) string {
	data := make([]byte, 512)
	fileStream.Read(data)

	t := http.DetectContentType(data)
	fileStream.Seek(0, io.SeekStart)
	return t
}

func getContentTypeBytes(fileStream *bytes.Buffer) string {
	data := make([]byte, 512)
	fileStream.Read(data)

	t := http.DetectContentType(data)
	fileStream.Reset()
	return t
}

func (s3s *S3Storage) ListAll() ([]string, error) {
	svc := s3.New(s3s.session)
	token := (*string)(nil)
	ret := make([]string, 0)
	for {
		objs, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:            aws.String(s3s.bucket),
			MaxKeys:           aws.Int64(1000),
			ContinuationToken: token,
		})
		token = objs.NextContinuationToken
		if err != nil {
			return nil, err
		}
		for i := range objs.Contents {
			ret = append(ret, *objs.Contents[i].Key)
			//fmt.Printf("Key:%s, Size:%d, ETag:%s, PartNumber:%d, StorageClass:%v\n",
			//	v.Contents[i].Key, v.Contents[i].Size, v.Contents[i].ETag, v.Contents[i].PartNumber, v.Contents[i].StorageClass)
		}
		if !*objs.IsTruncated {
			break
		}
	}

	return ret, nil
}
func (s *S3Storage) UploadFile(ctx context.Context, filePath string, fileStream multipart.File) error {
	uploader := s3manager.NewUploader(s.session)
	//contentType := getContentType(fileStream)

	extension, err := s.fetchFileContentType(ctx, filePath)
	if err != nil {
		fmt.Println("Fetch extension failed, err: ", err, ", key:", filePath)
		return err
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(filePath),
		Body:        fileStream,
		ContentType: aws.String(extension),
	})
	if err != nil {
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
		return "application/msword", nil
	case "pptx":
		return "application/vnd.ms-powerpoint", nil
	case "doc":
		return "application/msword", nil
	case "ppt":
		return "application/vnd.ms-powerpoint", nil
	case "xls":
		return "application/vnd.ms-excel", nil
	case "xlsx":
		return "application/vnd.ms-excel", nil
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
		region:     c.Region,
		endpoint:   c.Endpoint,
		accelerate: c.Accelerate,
		secretID: 	c.SecretID,
		secretKey: 	c.SecretKey,
	}
}

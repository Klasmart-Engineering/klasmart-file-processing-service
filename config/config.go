package config

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Config struct {
	Storage StorageConfig
	Log     LogConfig
	Core    CoreConfig
}

type CoreConfig struct {
	ExifToolPath string
	EyeD3Path    string
	//comma separated
	Processors string
}

type LogConfig struct {
	Level      string
	FailedFile string
	StdOut     bool
}

type StorageConfig struct {
	Driver string

	Accelerate bool
	EndPoint   string
	Bucket     string
	BucketOut  string
	Region     string
	AWSSession *session.Session
}

var (
	cfg Config
)

func LoadEnv(session *session.Session) {
	cfg = Config{
		Storage: StorageConfig{
			Driver:     getEnvStr("storage_driver", cfg.Storage.Driver),
			Accelerate: getEnvBool("storage_accelerate", cfg.Storage.Accelerate),
			EndPoint:   getEnvStr("storage_endpoint", cfg.Storage.EndPoint),
			Bucket:     getEnvStr("storage_bucket", cfg.Storage.Bucket),
			BucketOut:  getEnvStr("storage_bucket_out", cfg.Storage.Bucket),
			Region:     getEnvStr("storage_region", cfg.Storage.Region),
			AWSSession: session,
		},
		Log: LogConfig{
			Level:      getEnvStr("log_level", cfg.Log.Level),
			FailedFile: getEnvStr("log_failed_file", cfg.Log.FailedFile),
			StdOut:     getEnvBool("log_std_out", cfg.Log.StdOut),
		},
		Core: CoreConfig{
			ExifToolPath: getEnvStr("core__exiftool", cfg.Core.ExifToolPath),
			EyeD3Path:    getEnvStr("core__eyed3", cfg.Core.EyeD3Path),
			Processors:   getEnvStr("processors", cfg.Core.Processors),
		},
	}
}
func MustLoad(session *session.Session) {
	LoadEnv(session)
}

func Get() Config {
	return cfg
}

func GetRegion() (string, error) {
	var region = getEnvStr("storage_region", cfg.Storage.Region)
	if region == "" {
		return "", errors.New("storage_region environment variable must be set")
	}
	return region, nil
}

package config

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Storage StorageConfig `yaml:"storage"`
	API     APIConfig     `yaml:"api"`
	Log     LogConfig     `yaml:"log"`
	Core    CoreConfig    `yaml:"core"`
}

type CoreConfig struct {
	ExifToolPath string `yaml:"exiftool_path"`
	EyeD3Path    string `yaml:"eyeD3_path"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	FailedFile string `yaml:"failed_file"`
	StdOut     bool   `yaml:"std_out"`
}

type APIConfig struct {
	Port      int    `yaml:"port"`
	SecretKey string `yaml:"secret_key"`
}

type StorageConfig struct {
	Driver string `yaml:"driver"`

	Accelerate bool   `yaml:"accelerate"`
	EndPoint   string `yaml:"end_point"`
	Bucket     string `yaml:"bucket"`
	BucketOut  string `yaml:"bucket_out"`
	Region     string `yaml:"region"`
	AWSSession *session.Session
	//SecretID  string `yaml:"secret_id"`
	//SecretKey string `yaml:"secret_key"`
}

var (
	cfg Config
)

func LoadYAML(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &cfg)
	return err
}

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
			//SecretID:   getEnvStr("AWS_ACCESS_KEY_ID", cfg.Storage.SecretID),
			//SecretKey:  getEnvStr("AWS_SECRET_ACCESS_KEY", cfg.Storage.SecretKey),
		},
		API: APIConfig{
			Port:      getEnvInt("api_port", cfg.API.Port),
			SecretKey: getEnvStr("api_secret_key", cfg.API.SecretKey),
		},
		Log: LogConfig{
			Level:      getEnvStr("log_level", cfg.Log.Level),
			FailedFile: getEnvStr("log_failed_file", cfg.Log.FailedFile),
			StdOut:     getEnvBool("log_std_out", cfg.Log.StdOut),
		},
		Core: CoreConfig{
			ExifToolPath: getEnvStr("core__exiftool", cfg.Core.ExifToolPath),
			EyeD3Path:    getEnvStr("core__eyed3", cfg.Core.EyeD3Path),
		},
	}
}
func MustLoad(yaml string, session *session.Session) {
	LoadYAML(yaml)
	LoadEnv(session)
}

func Get() Config {
	return cfg
}

package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Storage StorageConfig `yaml:"storage"`
	MQ      MQConfig      `yaml:"mq"`
	API     APIConfig     `yaml:"api"`
	Log     LogConfig     `yaml:"log"`
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

type MQConfig struct {
	Driver                 string `yaml:"driver"`
	RedisHost              string `yaml:"redis_host"`
	RedisPort              int    `yaml:"redis_port"`
	RedisPassword          string `yaml:"redis_password"`
	RedisFailedPersistence string `yaml:"redis_failed_persistence"`

	MaxWorker int `yaml:"max_worker"`
}
type StorageConfig struct {
	Driver string `yaml:"driver"`

	Accelerate      bool   `yaml:"accelerate"`
	StorageEndPoint string `yaml:"storage_end_point"`
	StorageBucket   string `yaml:"bucket"`
	StorageRegion   string `yaml:"region"`

	SecretID  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}

var (
	cfg Config
)

func LoadYAML(path string) error{
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &cfg)
	return err
}

func LoadEnv() {
	cfg = Config{
		Storage: StorageConfig{
			Driver:          getEnvStr("storage.driver", cfg.Storage.Driver),
			Accelerate:      getEnvBool("storage.accelerate", cfg.Storage.Accelerate),
			StorageEndPoint: getEnvStr("storage.endpoint", cfg.Storage.StorageEndPoint),
			StorageBucket:   getEnvStr("storage.bucket", cfg.Storage.StorageBucket),
			StorageRegion:   getEnvStr("storage.region", cfg.Storage.StorageRegion),
			SecretID:        getEnvStr("storage.secret_id", cfg.Storage.SecretID),
			SecretKey:       getEnvStr("storage.secret_key", cfg.Storage.SecretKey),
		},
		MQ:      MQConfig{
			Driver:                 getEnvStr("mq.driver", cfg.MQ.Driver),
			RedisHost:              getEnvStr("mq.redis_host", cfg.MQ.RedisHost),
			RedisPort:              getEnvInt("mq.redis_port", cfg.MQ.RedisPort),
			RedisPassword:          getEnvStr("mq.redis_password", cfg.MQ.RedisPassword),
			RedisFailedPersistence: getEnvStr("mq.redis_failed_persistence", cfg.MQ.RedisFailedPersistence),
			MaxWorker:              getEnvInt("mq.max_worker", cfg.MQ.MaxWorker),
		},
		API:     APIConfig{
			Port:      getEnvInt("api.port", cfg.API.Port),
			SecretKey: getEnvStr("api.secret_key", cfg.API.SecretKey),
		},
		Log:     LogConfig{
			Level:      getEnvStr("log.level", cfg.Log.Level),
			FailedFile: getEnvStr("log.failed_file", cfg.Log.FailedFile),
			StdOut:     getEnvBool("log.std_out", cfg.Log.StdOut),
		},
	}
}
func MustLoad(yaml string)  {
	LoadYAML(yaml)
	LoadEnv()
}

func Get() Config {
	return cfg
}

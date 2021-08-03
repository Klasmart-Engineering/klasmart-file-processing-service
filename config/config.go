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
	cfg *Config
)

func LoadYAML(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	cfg = new(Config)
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		panic(err)
	}
}

func Get() *Config {
	if cfg == nil {
		panic("Config is nil")
	}
	return cfg
}

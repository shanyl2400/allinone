package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	GomssPath  string `yaml:"gomss_path"`
	BoltDBPath string `yaml:"boltdb_path"`

	HttpPort int `yaml:"http_port"`
}

var (
	_config     *Config
	_configOnce sync.Once
)

func GetConfig() *Config {
	_configOnce.Do(func() {
		_config = &Config{
			GomssPath:  "/home/shanyonglong/gomss",
			HttpPort:   8088,
			BoltDBPath: "./data",
		}
		conf, err := os.ReadFile("./config.yaml")
		if err != nil {
			log.Println("read config file failed")
			return
		}
		err = yaml.Unmarshal(conf, _config)
		if err != nil {
			log.Println("bad config file format")
			return
		}
	})
	return _config
}

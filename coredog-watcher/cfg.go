package coredogwatcher

import (
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	cfg              *Config
	onceCfg          sync.Once
	DEFAULT_CFG_PATH = "/etc/config/watcher.yaml"
)

const (
	COREFILE_DIR = "/corefile"
)

type Config struct {
	StorageConfig struct {
		Enabled bool `yaml:"enabled" env-default:"true"`
		// default storage protocol: s3
		Protocol          string `yaml:"protocol" env-default:"s3s"`
		S3AccessKeyID     string `yaml:"s3AccesskeyID"`
		S3SecretAccessKey string `yaml:"s3SecretAccessKey"`
		S3Region          string `yaml:"s3Region"`
		S3Bucket          string `yaml:"S3Bucket"`
		S3Endpoint        string `yaml:"S3Endpoint"`
		StoreDir          string `yaml:"StoreDir"`
		// presigned url expire time(by seconds)
		PresignedURLExpireSeconds int  `yaml:"PresignedURLExpireSeconds"`
		DeleteLocalCorefile       bool `yaml:"deleteLocalCorefile"`
	} `yaml:"StorageConfig"`
	Gc           bool   `yaml:"gc" env-default:"false"`
	GcType       string `yaml:"gc_type" env-default:"rm"`
	COREFILE_DIR string `yaml:"CorefileDir"`
}

func getCfg() *Config {
	onceCfg.Do(func() {
		cfg = &Config{}
		cfgPath := os.Getenv("CONFIG_WATCHER_PATH")
		if cfgPath == "" {
			cfgPath = DEFAULT_CFG_PATH
		}
		if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
			log.Fatal(err)
		}
		if cfg.COREFILE_DIR == "" {
			cfg.COREFILE_DIR = COREFILE_DIR
		}
	})
	return cfg
}

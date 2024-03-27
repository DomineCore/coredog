package coredogcontroller

import (
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	cfg              *Config
	onceCfg          sync.Once
	DEFAULT_CFG_PATH = "/etc/config/controller.yaml"
)

type Config struct {
	NoticeChannel []struct {
		// supported: wechat，slack
		Chan       string `yaml:"chan"`
		Webhookurl string `yaml:"webhookurl"`
	} `yaml:"NoticeChannel"`
}

func getCfg() *Config {
	onceCfg.Do(func() {
		cfg = &Config{}
		cfgPath := os.Getenv("CONFIG_CONTROLLER_PATH")
		if cfgPath == "" {
			cfgPath = DEFAULT_CFG_PATH
		}
		if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
			log.Fatal(err)
		}
	})
	return cfg
}

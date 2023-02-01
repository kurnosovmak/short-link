package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/kurnosovmak/short-link/pkg/logging"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	}
	Redis struct {
		Addr     string `yaml:"addr" env-default:"localhost:6379"`
		Password string `yaml:"password" env-default:""`
		DB       int    `yaml:"db" env-default:"0"`
	}
}

var instance *Config
var once sync.Once

func Get() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}

func (cfg *Config) GetFullAddress() string {
	return cfg.Listen.BindIP + ":" + cfg.Listen.Port
}

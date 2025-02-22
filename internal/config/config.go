package config

import (
	"errors"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	Env string `yaml:"env"`

	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`

	DB struct {
		FilePath    string        `yaml:"file_path"`
		OpenTimeout time.Duration `yaml:"open_timeout"`
	} `yaml:"db"`

	Yookassa struct {
		CheckoutURL string `env:"YOOKASSA_CHECKOUT_URL"`
		ShopID      string `env:"YOOKASSA_SHOP_ID"`
		APIToken    string `env:"YOOKASSA_API_TOKEN"`
	}
}

func GetConfig(path string) (*Config, error) {
	var err error

	once.Do(func() {
		instance = &Config{}

		err = errors.Join(
			cleanenv.ReadConfig(path, instance),
			godotenv.Load(),
			cleanenv.ReadEnv(instance),
		)
	})

	return instance, err
}

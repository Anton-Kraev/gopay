package bot

import (
	"errors"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	Token     string `env:"TG_BOT_TOKEN" env-required:""`
	AdminIDs  string `env:"TG_ADMIN_IDS" env-required:""`
	ServerURL string `env:"SERVER_URL" env-required:""`
}

func GetConfig() (*Config, error) {
	var err error

	once.Do(func() {
		instance = &Config{}

		err = errors.Join(
			godotenv.Load(),
			cleanenv.ReadEnv(instance),
		)
	})

	return instance, err
}

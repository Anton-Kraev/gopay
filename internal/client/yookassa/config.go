package yookassa

type Config struct {
	ID    string `env:"YOOKASSA_ID"`
	Token string `env:"YOOKASSA_TOKEN"`
	URL   string `env:"YOOKASSA_URL"`
}

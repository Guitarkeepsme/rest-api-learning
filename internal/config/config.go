package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	Clients     ClientsConfig `yaml:"clients"`
	AppSecret   string        `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
}

type HTTPServer struct {
	Addr        string        `yaml:"addr" env-default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	Username    string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoadConfig() *Config {
	// Надо брать из переменной окружения, а адрес должен быть относительным. Но почему-то
	// у меня ни то, ни другое на работает. Потом исправлю, а пока вот так:
	err := os.Setenv("CONFIG_PATH", "/Users/alexeykashurnikov/golang/rest api learning/url-shortener/config/local.yaml")
	if err != nil {
		log.Fatal(err)
	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH env var is not set")
	}

	// check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
	Insecure     bool          `yaml:"insecure"`
}

type ClientsConfig struct {
	SSO Client `yaml:"sso"`
}

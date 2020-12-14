package aocweb

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Config struct {
	GitHub     GitHub
	SessionKey string `env:"AOC_SESSION_KEY"`
	DB         DB
}

type GitHub struct {
	ClientID     string `env:"AOC_GH_CLIENT_ID"`
	ClientSecret string `env:"AOC_GH_CLIENT_SECRET"`

	OAuth2 *oauth2.Config
}

type DB struct {
	Host         string `env:"DB_HOST" envDefault:"127.0.0.1"`
	Port         string `env:"DB_PORT" envDefault:"5432"`
	Username     string `env:"DB_USER" envDefault:"postgres"`
	Password     string `env:"DB_PASSWORD,required"`
	DBName       string `env:"DB_NAME" envDefault:"aocweb"`
	SSLMode      string `env:"DB_SSL_MODE" envDefault:"disable"`
	MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"20"`
	MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"20"`
	Debug        bool   `env:"DB_DEBUG"`
}

func (c DB) DSN() string {
	dsn := "host=%s port=%s user=%s password='%s' dbname=%s sslmode=%s"
	return fmt.Sprintf(dsn,
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
}

func MustConfig() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	cfg.GitHub.OAuth2 = &oauth2.Config{
		ClientID:     cfg.GitHub.ClientID,
		ClientSecret: cfg.GitHub.ClientSecret,
		Endpoint:     github.Endpoint,
	}

	return &cfg
}

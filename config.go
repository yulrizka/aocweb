package aocweb

import (
	"log"

	"github.com/caarlos0/env/v6"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Config struct {
	GitHub     GitHub
	SessionKey string `env:"AOC_SESSION_KEY"`
}

type GitHub struct {
	ClientID     string `env:"AOC_GH_CLIENT_ID"`
	ClientSecret string `env:"AOC_GH_CLIENT_SECRET"`

	OAuth2 *oauth2.Config
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

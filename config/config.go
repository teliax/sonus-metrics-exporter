package config

import (
	"os"
	"strings"

	cfg "github.com/infinityworks/go-common/config"
)

// Config struct holds all the runtime configuration for the application
type Config struct {
	*cfg.BaseConfig
	APIURLs            []string
	APIUser            string
	APIPass            string
	APIAddressContexts []string
}

// Init populates the Config struct based on environmental runtime configuration
func Init() Config {

	c := cfg.Init()
	rawURLs := cfg.GetEnv("API_URLS", "https://172.16.7.2/api")
	rawACs := cfg.GetEnv("API_ADDRESSCONTEXTS", "default")
	user := os.Getenv("API_USER")
	pass := os.Getenv("API_PASSWORD")

	addressContexts := strings.Fields(rawACs)
	urls := strings.Fields(rawURLs)

	appConfig := Config{
		&c,
		urls,
		user,
		pass,
		addressContexts,
	}

	return appConfig
}

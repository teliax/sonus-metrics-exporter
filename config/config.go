package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	cfg "github.com/infinityworks/go-common/config"
)

// Config struct holds all the runtime configuration for the application
type Config struct {
	*cfg.BaseConfig
	APIURLs            []string
	APIUser            string
	APIPass            string
	APIAddressContexts []string
	APITimeout         time.Duration
}

// Init populates the Config struct based on environmental runtime configuration
func Init() Config {

	c := cfg.Init()
	rawURLs := cfg.GetEnv("API_URLS", "https://172.16.7.2/api")
	rawACs := cfg.GetEnv("API_ADDRESSCONTEXTS", "default")
	rawTimeout := cfg.GetEnv("API_TIMEOUT", "10")

	user := os.Getenv("API_USER")
	pass := os.Getenv("API_PASSWORD")

	addressContexts := strings.Fields(rawACs)
	urls := strings.Fields(rawURLs)

	intTimeout, err := strconv.ParseInt(rawTimeout, 0, 0)
	timeout := time.Duration(intTimeout) * time.Second

	if err != nil {
		panic("Unable to parse API_TIMEOUT as integer.")
	}

	appConfig := Config{
		&c,
		urls,
		user,
		pass,
		addressContexts,
		timeout,
	}

	return appConfig
}

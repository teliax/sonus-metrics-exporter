package config

import (
	"fmt"
	"os"

	cfg "github.com/infinityworks/go-common/config"
	log "github.com/sirupsen/logrus"
)

// Config struct holds all of the runtime configuration for the application
type Config struct {
	*cfg.BaseConfig
	APIURL     string
	APIUser    string
	APIPass    string
	TargetURLs []string
}

// Init populates the Config struct based on environmental runtime configuration
func Init() Config {

	ac := cfg.Init()
	url := cfg.GetEnv("API_URL", "https://172.16.7.2/api")
	user := os.Getenv("API_USER")
	pass := os.Getenv("API_PASSWORD")
	scraped, err := getScrapeURLs(url, user, pass)

	if err != nil {
		log.Errorf("Error initialising Configuration, Error: %v", err)
	}

	appConfig := Config{
		&ac,
		url,
		user,
		pass,
		scraped,
	}

	return appConfig
}

// Init populates the Config struct based on environmental runtime configuration
// All URL's are added to the TargetURL's string array
func getScrapeURLs(apiURL string, apiUser string, apiPass string) ([]string, error) {

	urls := []string{}

	opts := "" // "?&per_page=100" // Used to set extra API options.

	url := fmt.Sprintf("%s/operational/global/globalTrunkGroupStatus/%s", apiURL, opts)
	urls = append(urls, url)

	return urls, nil
}

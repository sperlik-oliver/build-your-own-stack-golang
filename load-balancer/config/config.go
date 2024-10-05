package config

import (
	"build-your-own/load-balancer/pool"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

var config Config

type Config struct {
	Upstreams []struct {
		URL string `json:"URL"`
	} `json:"upstreams"`
	HealthcheckIntervalSeconds int    `json:"healthcheckIntervalSeconds"`
	RetryIntervalSeconds       int    `json:"retryIntervalSeconds"`
	Port                       int    `json:"port"`
	Mode                       string `json:"mode"`
}

// Loads configuration from config.json
func LoadConfig() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Failed loading config from config.json: error opening file")
	}
	configByteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal("Failed loading config from config.json: error reading file")
	}
	unmarshalErr := json.Unmarshal(configByteValue, &config)
	if unmarshalErr != nil {
		log.Fatal("Failed loading config from config.json: error reading JSON")
	}

	if len(config.Upstreams) < 1 {
		log.Fatal("No upstreams defined in config")
	}

	for i := 0; i < len(config.Upstreams); i++ {
		configURL := config.Upstreams[i].URL
		parsedURL, err := url.Parse(configURL)

		if err != nil {
			log.Fatal("Failed loading config from config.json: error parsing upstream URL")
		}

		pool.AddUpstream(parsedURL)
	}

	if config.Mode != "roundrobin" && config.Mode != "leastconnections" {
		log.Fatalf("Unsupported load balancing mode: %s. Supported modes: 'roundrobin', 'leastconnections'", config.Mode)
	}

	log.Printf("Successfully loaded config from config.json: \n%+v", &config)
}

func GetConfig() *Config {
	return &config
}

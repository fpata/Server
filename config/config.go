package config

import (
	"clinic_server/cache"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port      string
	ServerUrl string
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	ConnectionString string
}

func GetConfiguration() Configurations {
	var config Configurations
	var isError bool = false
	if value, found := cache.CacheInstance.Get("Configurations"); found {
		if res, ok := value.(Configurations); ok {
			config = res
		}

	} else {

		jsonFile, err := os.Open("config.json")
		if err != nil {
			log.Error().Err(err).Msg("Failed to open config file")
			isError = true
		}
		defer jsonFile.Close()

		// Read the JSON file
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read config file")
			isError = true
		}

		// Unmarshal the JSON data into the Config struct

		if err := json.Unmarshal(byteValue, &config); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal JSON")
			isError = true
		}
		if !isError {
			cache.CacheInstance.Set("Configurations", config, time.Duration(time.Now().Year()))
		}
	}
	return config
}

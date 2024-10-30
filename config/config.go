package config

import (
	"clinic_server/cache"
	"clinic_server/logger"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
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
	logger.Init(zerolog.DebugLevel)
	var config Configurations
	var isError bool = false
	if value, found := cache.CacheInstance.Get("Configurations"); found {
		if res, ok := value.(Configurations); ok {
			config = res
		}

	} else {

		jsonFile, err := os.Open("config.json")
		if err != nil {

			logger.Error("Failed to open config file", err)
			isError = true
		}
		defer jsonFile.Close()

		// Read the JSON file
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			logger.Error("Failed to read config file", err)
			isError = true
		}

		// Unmarshal the JSON data into the Config struct

		if err := json.Unmarshal(byteValue, &config); err != nil {
			logger.Error("Failed to unmarshal JSON", err)
			isError = true
		}
		if !isError {
			cache.CacheInstance.Set("Configurations", config, time.Duration(time.Now().Year()))
		}
	}
	return config
}

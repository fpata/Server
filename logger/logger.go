package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger is the global logger instance
var Logger zerolog.Logger

// once ensures the logger is initialized only once
var once sync.Once

// Init initializes the logger with the specified log level
func Init(logLevel zerolog.Level) {
	once.Do(func() {
		// Create or open the log file with the current date
		logFile := fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02"))
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open log file")
		}

		// Configure zerolog to write to the log file
		Logger = zerolog.New(file).With().Timestamp().Logger()

		// Set the global log level
		zerolog.SetGlobalLevel(logLevel)

		// Set the global time format for zerolog
		zerolog.TimeFieldFormat = time.RFC3339
	})
}

// Info logs an info message
func Info(msg string, fields ...interface{}) {
	Logger.Info().Fields(fields).Msg(msg)
}

// Warn logs a warning message
func Warn(msg string, fields ...interface{}) {
	Logger.Warn().Fields(fields).Msg(msg)
}

// Error logs an error message
func Error(msg string, fields ...interface{}) {
	Logger.Error().Fields(fields).Msg(msg)
}

// Debug logs a debug message
func Debug(msg string, fields ...interface{}) {
	Logger.Debug().Fields(fields).Msg(msg)
}

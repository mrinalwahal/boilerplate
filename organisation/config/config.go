package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Base configuration.
type config struct {
	Env            *string         `mapstructure:"env"`
	Port           *string         `mapstructure:"port"`
	Database       *Database       `mapstructure:"database"`
	Authentication *Authentication `mapstructure:"authentication"`
	Logs           *Logs           `mapstructure:"logs"`
}

// Set default values for configuration.
func (c *config) SetDefaults() {
	if c.Env == nil {
		c.Env = new(string)
		*c.Env = "dev"
	}
	if c.Port == nil {
		c.Port = new(string)
		*c.Port = "8080"
	}
	if c.Database == nil {
		c.Database = &Database{
			Engine: "sqlite",
			DSN:    "file::memory:?cache=shared",
		}
	}
}

// Database configuration.
type Database struct {
	Engine string `mapstructure:"engine"`
	DSN    string `mapstructure:"dsn"` // Data Source Name
}

// Get gorm dialector.
func (d *Database) Dialector() gorm.Dialector {
	switch d.Engine {
	case "postgres":
		return postgres.Open(d.DSN)
	default:
		return sqlite.Open(d.DSN)
	}
}

// Authentication configuration.
type Authentication struct {
	Method string `mapstructure:"method"`
	Key    struct {
		Algorithm string `mapstructure:"algorithm"`
		Key       string `mapstructure:"key"`
	} `mapstructure:"key"`
}

// Logs configuration.
type Logs struct {
	Engine  string `mapstructure:"engine"`
	Address string `mapstructure:"address"`
	Level   Level  `mapstructure:"level"`
}

type Level string

func (l Level) ToSlogLevel() slog.Level {
	switch l {
	case DEBUG:
		return slog.LevelDebug
	case INFO:
		return slog.LevelInfo
	case WARN:
		return slog.LevelWarn
	case ERROR:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

const (
	DEBUG Level = "debug"
	INFO  Level = "info"
	WARN  Level = "warn"
	ERROR Level = "error"
)

func Get() *config {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("unable to read config file, %v", err))
	}

	var c config
	err := viper.Unmarshal(&c)
	if err != nil {
		panic(fmt.Sprintf("unable to decode into struct, %v", err))
	}

	// Set default values.
	c.SetDefaults()

	return &c
}

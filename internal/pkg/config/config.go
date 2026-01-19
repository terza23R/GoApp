package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Templates TemplatesConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	DSN string
}

type TemplatesConfig struct {
	Path string
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)

	v.SetDefault("templates.path", "templates/*.html")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = v.ReadInConfig()

	cfg := &Config{
		Server: ServerConfig{
			Host: v.GetString("server.host"),
			Port: v.GetInt("server.port"),
		},
		Database: DatabaseConfig{
			DSN: v.GetString("database.dsn"),
		},
		Templates: TemplatesConfig{
			Path: v.GetString("templates.path"),
		},
	}

	if cfg.Database.DSN == "" {
		return nil, fmt.Errorf("database.dsn is empty (set in config.yaml or env DATABASE_DSN)")
	}

	return cfg, nil
}

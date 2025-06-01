package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"runtime"
	"time"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Postgres PostgresConfig `mapstructure:"postgres"`
}

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	RunMode string `mapstructure:"run_mode"`
}

type PostgresConfig struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	Username    string        `mapstructure:"username"`
	Password    string        `mapstructure:"password"`
	Database    string        `mapstructure:"database"`
	SSLMode     string        `mapstructure:"ssl_mode"`
	MaxConns    int32         `mapstructure:"max_conns"`
	MinConns    int32         `mapstructure:"min_conns"`
	MaxLifetime time.Duration `mapstructure:"max_lifetime"`
	MaxIdleTime time.Duration `mapstructure:"max_idle_time"`
	HealthCheck time.Duration `mapstructure:"health_check"`
}

func parseConfig(v *viper.Viper) (*Config, error) {
	var config Config
	err := v.Unmarshal(&config)
	if err != nil {
		log.Printf("unable to decode struct, %v\n", err)
	}
	return &config, nil
}

func localConfig(fileName string, fileType string) (*viper.Viper, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(fileName)
	v.AddConfigPath(filepath.Join(basePath, "config_files"))
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		var fileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &fileNotFoundError) {
			return nil, errors.New("file not found")
		}
	}

	return v, nil
}

func getConfigPath(env string) string {
	if env == "production" {
		return "config-production"
	}
	if env == "staging" {
		return "config-staging"
	}
	if env == "local" {
		return "config-local"
	}
	if env == "test" {
		return "config-test"
	}

	panic("no valid env")
}

func GetConfig(env string) (*Config, error) {
	configPath := getConfigPath(env)
	loader, err := localConfig(configPath, "yaml")
	if err != nil {
		log.Printf("unable to finc config, %v\n", err)
		errors.New("no config found")
	}
	return parseConfig(loader)
}

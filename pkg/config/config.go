package config

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper

	database *dbCreds
}

func (c *Config) Database() *dbCreds {
	return c.database
}

func NewConfig() *Config {
	return &Config{
		viper: viper.New(),
	}
}

type dbCreds struct {
	Engine       string `json:"engine"`
	Host         string `json:"host"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	Port         int    `json:"port"`
	DatabaseName string `json:"databaseName"`
}

func (c *Config) setDatabase(ctx context.Context) error {
	return nil
}

func (c *Config) Load(ctx context.Context) error {
	envVarsToKey := map[string]struct {
		required bool
		key      string
	}{
		"DATABASE_SECRET": {true, "db.secret"},
	}
	for envVar, requiredKey := range envVarsToKey {
		if os.Getenv(envVar) == "" && requiredKey.required {
			return fmt.Errorf("missing required env var %s", envVar)
		}

		err := c.viper.BindEnv(requiredKey.key, envVar)
		if err != nil {
			return fmt.Errorf("bind env: %w", err)
		}
	}

	if err := c.setDatabase(ctx); err != nil {
		return fmt.Errorf("set database: %w", err)
	}

	return nil
}

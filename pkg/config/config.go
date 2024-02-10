package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	SecretCode string
	Env        string
}

type ServerConfig struct {
	Port int
}
type DatabaseConfig struct {
	DSN string
}
type JWTConfig struct {
	Secret string
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")

	viper.SetConfigName("app")

	viper.SetConfigType("env")

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg.Server.Port = viper.GetInt("PORT")
	cfg.Database.DSN = viper.GetString("DB_CONNECTIONSTRING")
	cfg.JWT.Secret = viper.GetString("JWT_ACCESS_SECRET")
	cfg.SecretCode = viper.GetString("SECRET_CODE")
	cfg.Env = viper.GetString("ENV")

	return &cfg, nil
}

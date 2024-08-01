package config

import (
	"hypha/api/internal/utils/logging"
	"os"

	"gopkg.in/yaml.v2"
)

var log = logging.Logger

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Dbname   string `yaml:"dbname"`
		Password string `yaml:"password"`
		Sslmode  string `yaml:"sslmode"`
	} `yaml:"database"`
	Http struct {
		Port       int `yaml:"port"`
		CorsPolicy struct {
			AllowOrigins     []string `yaml:"allow-origins"`
			AllowMethods     []string `yaml:"allow-methods"`
			AllowHeaders     []string `yaml:"allow-headers"`
			ExposeHeaders    []string `yaml:"expose-headers"`
			AllowCredentials bool     `yaml:"allow-credentials"`
			MaxAge           int      `yaml:"max-age"`
		} `yaml:"cors-policy"`
	} `yaml:"http"`
}

func ReadConfig(filename string) (*Config, error) {
	log.Info().Msgf("Reading configuration file %s", filename)

	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Error().Err(err).Msgf("Error reading file: %s", filename)
		return nil, err
	}
	log.Info().Msgf("Successfully read file: %s", filename)

	var cfg Config
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		log.Error().Err(err).Msgf("Error unmarshalling YAML from file: %s", filename)
		return nil, err
	}
	log.Info().Msgf("Successfully unmarshalled YAML from file: %s", filename)

	return &cfg, nil
}

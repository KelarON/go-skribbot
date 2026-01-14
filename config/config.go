package config

import (
	"fmt"
	"go-skribbot/model"
	"go-skribbot/utils"
	"log"
	"os"

	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

var ErrorNoConfigFile = errors.New("Config file not found")

// Config is the main configuration struct
type Config struct {
	PositionX   int               `yaml:"posotion_x"`
	PositionY   int               `yaml:"position_y"`
	DrawingType model.DrawingType `yaml:"drawing_type"`
}

// LoadConfig loads the config from config file
func LoadConfig() (*Config, error) {

	cfg := &Config{}

	if !configFileExists() {
		cfg.DrawingType = model.DRAWING_TYPE_LINE
		return cfg, ErrorNoConfigFile
	}

	err := cleanenv.ReadConfig("config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	return cfg, nil
}

// Save saves the config to config file
func (cfg *Config) Save() {

	err := utils.WriteStructToYAMLFile("config.yaml", cfg)
	if err != nil {
		log.Printf("error saving config file: %v", err)
		return
	}
}

// configFileExists checks if config file exists
func configFileExists() bool {
	_, err := os.Stat("config.yaml")
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

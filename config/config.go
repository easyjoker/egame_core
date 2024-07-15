package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var ConfigData *Config

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	MaxIdle  int    `yaml:"max_idle"`
	MaxOpen  int    `yaml:"max_open"`
	MaxLife  int    `yaml:"max_life"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}

type Config struct {
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
	Server   Server   `yaml:"server"`
}

func LoadConfig() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("ReadFile: %v", err)
		return
	}
	err = yaml.Unmarshal(data, ConfigData)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func GetConfig() *Config {
	if ConfigData == nil {
		LoadConfig()
	}
	return ConfigData
}

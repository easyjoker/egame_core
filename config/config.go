package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var ConfigData *Config
var resPath string

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
	DB       int    `yaml:"db"`
}

type Config struct {
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
	Server   Server   `yaml:"server"`
}

func loadConfig() {
	data, err := os.ReadFile(resPath)
	if err != nil {
		log.Fatalf("ReadFile: %v", err)
		return
	}
	var con Config
	err = yaml.Unmarshal(data, &con)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	ConfigData = &con
}

func Init(path string) {
	resPath = path
	loadConfig()
}

func GetConfig() *Config {
	if ConfigData == nil {
		loadConfig()
	}
	return ConfigData
}

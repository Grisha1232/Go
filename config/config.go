package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config структура для хранения параметров конфигурации
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// ServerConfig настройки сервера
type ServerConfig struct {
	Port int `yaml:"port"`
}

// DatabaseConfig настройки базы данных
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// JWTConfig настройки JWT-токенов
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	Expiration int    `yaml:"expiration"` // Время жизни токена в секундах
}

// LoadConfig загружает конфигурацию из файла `config.yaml`
func LoadConfig() (*Config, error) {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка открытия файла конфигурации: %v", err)
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Ошибка декодирования файла конфигурации: %v", err)
		return nil, err
	}

	return &config, nil
}

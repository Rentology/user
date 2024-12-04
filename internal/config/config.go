package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	RabbitMQ RabbitMQConfig `yaml:"rabbit_mq"`
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Env     string `yaml:"env"`
	Version string `yaml:"version"`
}

type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Ssl          bool   `yaml:"ssl"`
	JwtSecretKey string `yaml:"jwt_secret_key"`
	Csrf         bool   `yaml:"csrf"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
	SslMode  bool   `yaml:"sslMode"`
}

type RabbitMQConfig struct {
	Url       string `yaml:"url"`
	QueueName string `yaml:"queue_name"`
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s doesn't exists", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}
	return &cfg
}

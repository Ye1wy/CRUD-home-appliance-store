package config

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/database/connection"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServiceId                 string `env:"service_id" env-default:"go-service-0001"`
	ServiceName               string `env:"service_name" env-default:"crud-service"`
	Env                       string `env:"env" env-default:"local"`
	connection.PostgresConfig `env:"POSTGRES"`
	HTTPServer                `env:"http_server"`
}

type HTTPServer struct {
	Address string `env:"address" env-default:"localhost"`
	Port    string `env:"port" env-default:"8080"`
}

func MustLoad() *Config {
	op := "config.MustLoad"

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("op: %s, CONFIG_PATH is empty", op)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("op: %s, Error: There is no config file: %v", op, err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("op: %v, Error: Cannot read config: %v", op, err)
	}

	return &cfg
}

func (cfg *Config) PrintInfo() {
	fmt.Println("---------------------")
	fmt.Println("env: " + cfg.Env)
	fmt.Println("address: " + cfg.Address)
	fmt.Println("port: " + cfg.Port)
	fmt.Println("DNS: " + cfg.PostgresHost)
	fmt.Println("---------------------")
}

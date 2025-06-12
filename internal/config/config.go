package config

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/database/connection"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `env:"env" env-default:"local"`
	PostgresConfig connection.PostgresConfig
	CrudService    CrudService
	ConsulService  ConsulConfig
}

type CrudService struct {
	Id      string `env:"crud_service_id" env-default:"go-service-0001"`
	Name    string `env:"crud_service_name" env-default:"crud-service"`
	Address string `env:"crud_service_address" env-default:"localhost"`
	Port    string `env:"crud_service_port" env-default:"8080"`
}

type ConsulConfig struct {
	Address        string        `env:"consul_service_address"`
	Port           string        `env:"consul_service_port" env-default:"8500"`
	RetryDelay     time.Duration `env:"consul_service_retry_delay" env-default:"2s"`
	MaxAttempts    int           `env:"consul_service_max_attempts" env-default:"5"`
	SurveyInterval string        `env:"consul_service_survey_interval" env-default:"60s"`
	SurveyTimeout  string        `env:"consul_service_survey_timeout" env-default:"10s"`
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
	fmt.Println("address: " + cfg.CrudService.Address)
	fmt.Println("port: " + cfg.CrudService.Port)
	fmt.Println("---------------------")
}

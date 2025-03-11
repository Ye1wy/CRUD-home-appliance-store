package config

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `env:"env" env-default:"local"`
	MongoURL   string `env:"MONGO_URL" env-required:"true"`
	HTTPServer `env:"http_server"`
}

type HTTPServer struct {
	Address      string        `env:"address" env-default:"localhost"`
	Port         string        `env:"port" env-default:"8080"`
	Timeout      time.Duration `env:"timeout" env-default:"5s"`
	IddleTimeout time.Duration `env:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	op := "config.MustLoad"

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("[ERROR] CONFIG_PATH is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("[ERROR] There is no config file: %s", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("[ERROR] Cannot read config: ", logger.Err(err), "op", op)
	}

	return &cfg
}

func (cfg *Config) PrintInfo() {
	fmt.Println("---------------------")
	fmt.Println("env: " + cfg.Env)
	fmt.Println("address: " + cfg.Address)
	fmt.Println("port: " + cfg.Port)
	fmt.Println("mongourl: " + cfg.MongoURL)
	fmt.Println("---------------------")

}

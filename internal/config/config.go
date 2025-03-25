package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `env:"env" env-default:"local"`
	MongoURI   string `env:"MONGO_URI" env-required:"true"`
	HTTPServer `env:"http_server"`
}

type HTTPServer struct {
	Address string `env:"address" env-default:"localhost"`
	Port    string `env:"port" env-default:"8080"`
}

func MustLoad() *Config {
	op := "config.MustLoad"

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("Op: %s, CONFIG_PATH is empty", op)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Op: %s, Error: There is no config file: %v", op, err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Op: %v, Error: Cannot read config: %v", op, err)
	}

	return &cfg
}

func (cfg *Config) PrintInfo() {
	fmt.Println("---------------------")
	fmt.Println("env: " + cfg.Env)
	fmt.Println("address: " + cfg.Address)
	fmt.Println("port: " + cfg.Port)
	fmt.Println("mongourl: " + cfg.MongoURI)
	fmt.Println("---------------------")

}

package consul

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"fmt"
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

func Registration(cfg *config.Config) error {
	op := "consul.registration.Registration"

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = "http://consul-service:8500" // need to be moved to .env file
	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	reg := &consulapi.AgentServiceRegistration{
		ID:      cfg.ServiceId,
		Name:    cfg.ServiceName,
		Port:    port,
		Address: cfg.Address,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/api/check", cfg.Address, port),
			Interval: "10s",
			Timeout:  cfg.ConnectTimeout.String(),
		},
	}

	return client.Agent().ServiceRegister(reg)
}

func RetryRegistration(cfg *config.Config, log *logger.Logger) {
	const maxAttempts = 5
	const retryDelay = 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := Registration(cfg)
		if err != nil {
			log.Warn("Failed to register service in Consul, will retry...",
				"attempt", attempt,
				logger.Err(err),
			)
			time.Sleep(retryDelay)
			continue
		}

		log.Info("Service successfully registered in Consul")
		return
	}

	log.Error("Failed to register service in Consul after all retries")
}

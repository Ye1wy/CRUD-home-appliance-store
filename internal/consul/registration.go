package consul

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"fmt"
	"strconv"

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

package consul

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

func WaitForService(cfg *config.Config) error {
	consulConfig := &consulapi.Config{Address: fmt.Sprintf("http://%s:%s", cfg.ConsulService.Address, cfg.ConsulService.Port)}
	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		return fmt.Errorf("failed to create consul client: %w", err)
	}

	timeout, err := time.ParseDuration(cfg.ConsulService.SurveyTimeout)
	if err != nil {
		return fmt.Errorf("invalid time for calculating timeout: %w", err)
	}

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		services, _, err := client.Health().Service(cfg.CrudService.Name, "", true, nil)
		if err == nil && len(services) > 0 {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("timeout: service %s not available in consul", cfg.CrudService.Name)
}

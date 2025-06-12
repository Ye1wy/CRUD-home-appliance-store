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
	consulConfig.Address = fmt.Sprintf("http://%s:%s", cfg.ConsulService.Address, cfg.ConsulService.Port)
	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	servicePort, err := strconv.Atoi(cfg.CrudService.Port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	reg := &consulapi.AgentServiceRegistration{
		ID:      cfg.CrudService.Id,
		Name:    cfg.CrudService.Name,
		Port:    servicePort,
		Address: cfg.CrudService.Address,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/api/check", cfg.CrudService.Address, servicePort),
			Interval: cfg.ConsulService.SurveyInterval,
			Timeout:  cfg.ConsulService.SurveyTimeout,
		},
	}

	return client.Agent().ServiceRegister(reg)
}

func RetryRegistration(cfg *config.Config, log *logger.Logger) {
	for attempt := 1; attempt <= cfg.ConsulService.MaxAttempts; attempt++ {
		err := Registration(cfg)
		if err != nil {
			log.Warn("Failed to register service in Consul, will retry...",
				"attempt", attempt,
				logger.Err(err),
			)
			time.Sleep(cfg.ConsulService.RetryDelay)
			continue
		}

		log.Info("Service successfully registered in Consul")
		return
	}

	log.Error("Failed to register service in Consul after all retries")
}

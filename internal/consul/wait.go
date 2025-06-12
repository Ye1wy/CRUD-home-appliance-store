package consul

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

func WaitForService(cfg *config.Config) error {
	consulConfig := &consulapi.Config{Address: cfg.ConsulService.Address}
	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		return fmt.Errorf("failed to create consul client: %w", err)
	}

	services, _, err := client.Catalog().Services(nil)
	if err != nil {
		log.Fatal("Error retrieving services:", err)
	}

	fmt.Println("Registered Services:")
	for serviceName, tags := range services {
		fmt.Printf("- %s (Tags: %v)\n", serviceName, tags)
	}
	// timeout, err := time.ParseDuration(cfg.ConsulService.SurveyTimeout)
	// if err != nil {
	// 	return fmt.Errorf("invalid time for calculating timeout: %w", err)
	// }

	// deadline := time.Now().Add(timeout)

	// for time.Now().Before(deadline) {
	// 	services, _, err := client.Health().Service(cfg.CrudService.Name, "", true, nil)
	// 	if err == nil && len(services) > 0 {
	// 		return nil
	// 	}

	// 	time.Sleep(1 * time.Second)
	// }

	return fmt.Errorf("timeout: service %s not available in consul", cfg.CrudService.Name)
}

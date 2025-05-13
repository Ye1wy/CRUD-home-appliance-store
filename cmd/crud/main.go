package main

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/internal/controllers"
	"CRUD-HOME-APPLIANCE-STORE/internal/database/connection"
	repository "CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/internal/routes"
	"CRUD-HOME-APPLIANCE-STORE/internal/services"
	"CRUD-HOME-APPLIANCE-STORE/internal/uow"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	cfg := config.MustLoad()

	cfg.PrintInfo()
	log := logger.NewLogger(cfg.Env)
	log.Info("Logger is created")
	conn, err := connection.NewPostgresStorage(&cfg.PostgresConfig)
	if err != nil {
		log.Error("Error in connetion to postgres: ", logger.Err(err))
		os.Exit(1)
	}

	log.Info("Connection is established")
	unit := repository.NewUnitOfWork(conn, log)

	unit.Register("client", func(tx pgx.Tx, log *logger.Logger) uow.Repository {
		return postgres.NewClientRepository(tx, log)
	})

	clientRepo := postgres.NewClientRepository(conn, log)
	clientService := services.NewClientService(clientRepo, unit, log)
	clientController := controllers.NewClientsController(clientService, log)

	routerConfig := routes.RouterConfig{
		ClientController: clientController,
	}

	router := routes.NewRouter(routerConfig)
	log.Info("The paths are laid")

	log.Info("Server started")
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Error("Error in start server: ", logger.Err(err))
	}
}

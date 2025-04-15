package main

import (
	"CRUD-HOME-APPLIANCE-STORE/api/routes"
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/internal/controllers"
	"CRUD-HOME-APPLIANCE-STORE/internal/database/postgres"
	psgrep "CRUD-HOME-APPLIANCE-STORE/internal/repositories/postgres"
	"CRUD-HOME-APPLIANCE-STORE/internal/services"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"os"
)

func main() {
	cfg := config.MustLoad()

	cfg.PrintInfo()
	log := logger.NewLogger(cfg.Env)
	log.Info("Logger is created")
	conn, err := postgres.NewPostgresStorage(&cfg.PostgresConfig)
	if err != nil {
		log.Error("Error in connetion to postgres: ", logger.Err(err))
		os.Exit(1)
	}

	log.Info("Connection is established")

	clientRepo := psgrep.NewClientRepository(conn, log)
	clientService := services.NewClientService(clientRepo, clientRepo, log)
	clientController := controllers.NewClientsController(clientService, log)

	// productRepo := repositories.NewMongoProductsRepository(storage.Database, database.PRODUCTS, log)
	// productService := services.NewProductService(productRepo, log)
	// productController := controllers.NewProductController(productService, log)

	// supplierRepo := repositories.NewMongoSupplierRepository(storage.Database, database.SUPPLIERS, log)
	// supplierService := services.NewSupplierService(supplierRepo, log)
	// supplierController := controllers.NewSupplierContoller(supplierService, log)

	routerConfig := routes.RouterConfig{
		ClientController: clientController,
		// ProductController:  ,
		// SupplierController: ,
	}

	router := routes.NewRouter(routerConfig)
	log.Info("The paths are laid")

	log.Info("Server started")
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Error("Error in start server: ", logger.Err(err))
	}
}

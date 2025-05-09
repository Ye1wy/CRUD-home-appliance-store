package main

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/internal/database/connection"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"os"
)

func main() {
	cfg := config.MustLoad()

	cfg.PrintInfo()
	log := logger.NewLogger(cfg.Env)
	log.Info("Logger is created")
	_, err := connection.NewPostgresStorage(&cfg.PostgresConfig)
	if err != nil {
		log.Error("Error in connetion to postgres: ", logger.Err(err))
		os.Exit(1)
	}

	log.Info("Connection is established")

	// // clientRepo := postgres.NewClientRepository(conn, log)
	// // clientService := services.NewClientService(clientRepo, clientRepo, log)
	// // clientController := controllers.NewClientsController(clientService, log)

	// // productRepo := postgres.NewProductRepository(conn, log)
	// // productService := services.NewProductService(productRepo, productRepo, log)
	// // productController := controllers.NewProductController(productService, log)

	// // routerConfig := routes.RouterConfig{
	// // 	ClientController:  clientController,
	// // 	ProductController: productController,
	// // 	// SupplierController: ,
	// // }

	// // router := routes.NewRouter(routerConfig)
	// // log.Info("The paths are laid")

	// // log.Info("Server started")
	// // if err := router.Run(":" + cfg.Port); err != nil {
	// // 	log.Error("Error in start server: ", logger.Err(err))
	// }
}

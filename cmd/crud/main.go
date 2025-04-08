package main

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/internal/controllers"
	"CRUD-HOME-APPLIANCE-STORE/internal/database"
	"CRUD-HOME-APPLIANCE-STORE/internal/database/mongodb"
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"CRUD-HOME-APPLIANCE-STORE/internal/routes"
	"CRUD-HOME-APPLIANCE-STORE/internal/services"
	"os"
)

func main() {
	cfg := config.MustLoad()

	cfg.PrintInfo()
	cfg.PrintInfo()
	log := logger.NewLogger(cfg.Env)
	log.Info("Logger is created")
	storage, err := mongodb.Connect(cfg.MongoURI)
	log.Info("Connection is established")
	if err != nil {
		log.Error("Error in connetion to mongoDB: ", logger.Err(err))
		os.Exit(1)
	}

	clientRepo := repositories.NewMongoClientRepository(storage.Database, database.CLIENTS, log)
	clientService := services.NewClientService(clientRepo, log)
	clientController := controllers.NewClientsController(clientService, log)

	productRepo := repositories.NewMongoProductsRepository(storage.Database, database.PRODUCTS, log)
	productService := services.NewProductService(productRepo, log)
	productController := controllers.NewProductController(productService, log)

	supplierRepo := repositories.NewMongoSupplierRepository(storage.Database, database.SUPPLIERS, log)
	supplierService := services.NewSupplierService(supplierRepo, log)
	supplierController := controllers.NewSupplierContoller(supplierService, log)

	routerConfig := routes.RouterConfig{
		ClientController:   clientController,
		ProductController:  productController,
		SupplierController: supplierController,
	}

	router := routes.NewRouter(routerConfig)
	log.Info("The paths are laid")

	log.Info("Server started")
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Error("Error in start server: ", logger.Err(err))
	}
}

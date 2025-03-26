package main

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/internal/controllers"
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

	clientRepo := repositories.NewMongoClientRepository(storage.Database)
	log.Info("MongoDB client repository is created")
	clientService := services.NewClientService(clientRepo)
	log.Info("MongoDB client service is created")
	clientController := controllers.NewClientsController(clientService, log)
	log.Info("MongoDB client controller is created")

	productRepo := repositories.NewMongoProductsRepository(storage.Database)
	log.Info("MongoDB product repository is created")
	productService := services.NewProductService(productRepo)
	log.Info("MongoDB product service is created")
	productController := controllers.NewProductController(productService, log)
	log.Info("MongoDB product controller is created")

	routerConfig := routes.RouterConfig{
		ClientController:  clientController,
		ProductController: productController,
	}

	router := routes.NewRouter(routerConfig)
	log.Info("The paths are laid")

	log.Info("Server started")
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Error("Error in start server: ", logger.Err(err))
	}
}

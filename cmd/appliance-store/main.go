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
	log := logger.NewLogger(cfg.Env)
	storage, err := mongodb.Connect(cfg.MongoURL, log)
	if err != nil {
		log.Error("[Error] Error in connetion to mongoDB: ", logger.Err(err))
		os.Exit(1)
	}

	log.Info("[INFO] Server started")

	clientRepo := repositories.NewMongoClientRepository(storage.Database)
	clientService := services.NewClientService(clientRepo)
	clientController := controllers.NewClientsController(clientService, log)

	router := routes.NewRouter(clientController)
	router.Run(":8080")
}

package main

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/internal/database/mongodb"
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
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

	_ = storage

	log.Info("[INFO] Server started")

}

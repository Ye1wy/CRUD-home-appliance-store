package main

import (
	_ "CRUD-HOME-APPLIANCE-STORE/api"
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/internal/consul"
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

//	@title			Swagger CRUD Home appliance store API
//	@version		1.0
//	@description	This is a sample server Petstore server.

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		aboba.com
// @BasePath	/api/v1
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

	go func() {
		if err := consul.Registration(cfg); err != nil {
			log.Error("Failed to register the service in consul", logger.Err(err))
		}
	}()

	unit := repository.NewUnitOfWork(conn, log)

	err = unit.Register("client", func(tx pgx.Tx, log *logger.Logger) uow.Repository {
		return postgres.NewClientRepository(tx, log)
	})
	if err != nil {
		log.Error("Client repository registration in uow is unable")
		os.Exit(1)
	}

	err = unit.Register("address", func(tx pgx.Tx, log *logger.Logger) uow.Repository {
		return postgres.NewAddressRepository(tx, log)
	})
	if err != nil {
		log.Error("Address repository registration in uow is unable")
		os.Exit(1)
	}

	err = unit.Register("supplier", func(tx pgx.Tx, log *logger.Logger) uow.Repository {
		return postgres.NewSupplierRepository(tx, log)
	})
	if err != nil {
		log.Error("Supplier repository registration in uow is unable")
		os.Exit(1)
	}

	err = unit.Register("image", func(tx pgx.Tx, log *logger.Logger) uow.Repository {
		return postgres.NewImageRepository(tx, log)
	})
	if err != nil {
		log.Error("Image repository registration in uow is unable")
		os.Exit(1)
	}

	err = unit.Register("product", func(tx pgx.Tx, log *logger.Logger) uow.Repository {
		return postgres.NewProductRepository(tx, log)
	})
	if err != nil {
		log.Error("Product repository registration in uow is unable")
		os.Exit(1)
	}

	clientRepo := postgres.NewClientRepository(conn, log)
	clientService := services.NewClientService(clientRepo, unit, log)
	clientController := controllers.NewClientsController(clientService, log)

	supplierRepo := postgres.NewSupplierRepository(conn, log)
	supplierService := services.NewSupplierService(supplierRepo, unit, log)
	supplierController := controllers.NewSupplierContoller(supplierService, log)

	imageRepo := postgres.NewImageRepository(conn, log)
	imageService := services.NewImageService(imageRepo, unit, log)
	imageController := controllers.NewImageController(imageService, log)

	productRepo := postgres.NewProductRepository(conn, log)
	productService := services.NewProductService(productRepo, unit, log)
	productController := controllers.NewProductController(productService, log)

	routerConfig := routes.RouterConfig{
		ClientController:   clientController,
		ProductController:  productController,
		SupplierController: supplierController,
		ImageController:    imageController,
	}

	router := routes.NewRouter(routerConfig)
	log.Info("The paths are laid")

	log.Info("Server started")
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Error("Error in start server: ", logger.Err(err))
	}
}

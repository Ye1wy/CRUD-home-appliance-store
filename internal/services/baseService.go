package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
)

type BaseServiceInterface interface{}

type BaseService struct {
	Repo   *repositories.BaseRepositoryInterface
	Logger *logger.Logger
}

func NewBaseService(repository *repositories.BaseRepositoryInterface, logger *logger.Logger) *BaseService {
	logger.Debug("Base Service is created")
	return &BaseService{
		Repo:   repository,
		Logger: logger,
	}
}

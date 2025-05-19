package controllers

import (
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"

	"github.com/gin-gonic/gin"
)

type BaseControllerInterface interface {
}

type BaseController struct {
	logger *logger.Logger
}

func NewBaseContorller(logger *logger.Logger) *BaseController {
	return &BaseController{
		logger: logger,
	}
}

func (ctrl *BaseController) responce(c *gin.Context, statusCode int, obj any) {
	switch c.GetHeader("Accept") {
	case "appication/xml":
		c.XML(statusCode, obj)
	default:
		c.JSON(statusCode, obj)
	}
}

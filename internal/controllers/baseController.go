package controllers

import (
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"fmt"

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

func (ctrl *BaseController) mapping(c *gin.Context, obj any) error {
	switch c.GetHeader("content-type") {
	case "application/xml":
		if err := c.BindXML(&obj); err != nil {
			return fmt.Errorf("Failed to bind xml: %v", err)
		}
	default:
		if err := c.BindJSON(&obj); err != nil {
			return fmt.Errorf("Failed to bind json: %v", err)
		}
	}

	return nil
}

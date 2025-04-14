package routes

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type routes struct {
	router *gin.Engine
}

type RouterConfig struct {
	ClientController *controllers.ClientController
	// ProductController  *controllers.ProductController
	// SupplierController *controllers.SupplierController
}

func NewRouter(cfg RouterConfig) routes {
	r := routes{
		router: gin.Default(),
	}

	r.router.GET("/api/v1/swagger/index.html", swaggerGive)
	r.router.GET("/", func(c *gin.Context) {
		c.File("./misc/images/amogus.gif")
	})

	clientGroup := r.router.Group("/api/v1/clients")
	{
		clientGroup.GET("", cfg.ClientController.GetAll)
		clientGroup.POST("", cfg.ClientController.Create)
		clientGroup.GET("", cfg.ClientController.GetByNameAndSurname)
		clientGroup.PATCH("/:id", cfg.ClientController.UpdateAddress)
		clientGroup.DELETE("/:id", cfg.ClientController.Delete)
	}

	// productGroup := r.router.Group("/api/v1/products")
	// {
	// 	productGroup.GET("", cfg.ProductController.GetAll)
	// 	productGroup.POST("", cfg.ProductController.Create)
	// 	productGroup.GET("/:id", cfg.ProductController.GetById)
	// 	productGroup.PATCH("/:id/decrease", cfg.ProductController.DecreaseStock)
	// 	productGroup.DELETE("/:id", cfg.ProductController.Delete)
	// }

	// supplierGroup := r.router.Group("/api/v1/suppliers")
	// {
	// 	supplierGroup.GET("", cfg.SupplierController.GetAll)
	// 	supplierGroup.POST("", cfg.SupplierController.Create)
	// 	supplierGroup.GET("/:id", cfg.SupplierController.GetById)
	// 	supplierGroup.PATCH("/:id", cfg.SupplierController.UpdateAddress)
	// 	supplierGroup.DELETE("/:id", cfg.SupplierController.Delete)
	// }

	return r
}

func (r routes) Run(addr ...string) error {
	return r.router.Run(addr...)
}

func swaggerGive(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"massage": "Hello durik!"})
}

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
	ClientController  *controllers.ClientsController
	ProductController *controllers.ProductsController
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
		clientGroup.GET("", cfg.ClientController.GetAllClients)
		clientGroup.POST("", cfg.ClientController.AddClient)
		clientGroup.GET("/search", cfg.ClientController.GetClientByNameAndSurname)
		clientGroup.PATCH("/:id/address", cfg.ClientController.ChangeAddressParameter)
		clientGroup.DELETE("/:id", cfg.ClientController.DeleteClientById)
	}

	productGroup := r.router.Group("/api/v1/products")
	{
		productGroup.GET("", cfg.ProductController.GetAllProduct)
		productGroup.POST("", cfg.ProductController.AddProduct)
		productGroup.GET("/:id", cfg.ProductController.GetProductById)
		productGroup.PATCH("/:id/decrease", cfg.ProductController.DecreaseParameter)
		productGroup.DELETE("/:id", cfg.ProductController.DeleteProductById)
	}

	return r
}

func (r routes) Run(addr ...string) error {
	return r.router.Run(addr...)
}

func swaggerGive(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"massage": "Hello durik!"})
}

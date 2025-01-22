package routes

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type routes struct {
	router *gin.Engine
}

func NewRouter(clientController *controllers.ClientsController) routes {
	r := routes{
		router: gin.Default(),
	}

	r.router.GET("/api/v1/swagger/index.html", swaggerGive)

	clientGroup := r.router.Group("/api/v1/clients")
	{
		clientGroup.GET("", clientController.GetAllClients)
		clientGroup.POST("", clientController.AddClient)
		clientGroup.GET("/search", clientController.GetClientByNameAndSurname)
		clientGroup.DELETE("/:id", clientController.DeleteClientById)
	}

	return r
}

func (r routes) Run(addr ...string) error {
	return r.router.Run(addr...)
}

func swaggerGive(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"massage": "Hello durik!"})
}

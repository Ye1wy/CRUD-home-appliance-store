package routes

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// // Route is the information for every URI.
// type Route struct {
// 	// Name is the name of this Route.
// 	Name string
// 	// Method is the string for the HTTP method. ex) GET, POST etc..
// 	Method string
// 	// Pattern is the pattern of the URI.
// 	Pattern string
// 	// HandlerFunc is the handler function of this route.
// 	HandlerFunc gin.HandlerFunc
// }

// // NewRouter returns a new router.
// func NewRouter(handleFunctions ApiHandleFunctions) *gin.Engine {
// 	return NewRouterWithGinEngine(gin.Default(), handleFunctions)
// }

// // NewRouter add routes to existing gin engine.
// func NewRouterWithGinEngine(router *gin.Engine, handleFunctions ApiHandleFunctions) *gin.Engine {
// 	for _, route := range getRoutes(handleFunctions) {
// 		if route.HandlerFunc == nil {
// 			route.HandlerFunc = DefaultHandleFunc
// 		}

// 		switch route.Method {
// 		case http.MethodGet:
// 			router.GET(route.Pattern, route.HandlerFunc)
// 		case http.MethodPost:
// 			router.POST(route.Pattern, route.HandlerFunc)
// 		case http.MethodPut:
// 			router.PUT(route.Pattern, route.HandlerFunc)
// 		case http.MethodPatch:
// 			router.PATCH(route.Pattern, route.HandlerFunc)
// 		case http.MethodDelete:
// 			router.DELETE(route.Pattern, route.HandlerFunc)
// 		}
// 	}

// 	return router
// }

// // Default handler for not yet implemented routes
// func DefaultHandleFunc(c *gin.Context) {
// 	c.String(http.StatusNotImplemented, "501 not implemented")
// }

// type ApiHandleFunctions struct {
// 	// Routes for the ClientsAPI part of the API
// 	ClientsAPI ClientsAPI
// 	// Routes for the ImagesAPI part of the API
// 	ImagesAPI ImagesAPI
// 	// Routes for the ProductsAPI part of the API
// 	ProductsAPI ProductsAPI
// 	// Routes for the SuppliersAPI part of the API
// 	SuppliersAPI SuppliersAPI
// }

// func getRoutes(handleFunctions ApiHandleFunctions) []Route {
// 	return []Route{
// 		{
// 			"AddClient",
// 			http.MethodPost,
// 			"/api/v1/api/v1/clients",
// 			handleFunctions.ClientsAPI.AddClient,
// 		},
// 		{
// 			"ChangeAddressIdParameter",
// 			http.MethodPatch,
// 			"/api/v1/api/v1/clients/:id/address",
// 			handleFunctions.ClientsAPI.ChangeAddressIdParameter,
// 		},
// 		{
// 			"DeleteClientById",
// 			http.MethodDelete,
// 			"/api/v1/api/v1/clients/:id",
// 			handleFunctions.ClientsAPI.DeleteClientById,
// 		},
// 		{
// 			"GetAllClients",
// 			http.MethodGet,
// 			"/api/v1/api/v1/clients",
// 			handleFunctions.ClientsAPI.GetAllClients,
// 		},
// 		{
// 			"SearchClientByNameAndSurname",
// 			http.MethodGet,
// 			"/api/v1/api/v1/clients/search",
// 			handleFunctions.ClientsAPI.SearchClientByNameAndSurname,
// 		},
// 		{
// 			"AddImage",
// 			http.MethodPost,
// 			"/api/v1/api/v1/images",
// 			handleFunctions.ImagesAPI.AddImage,
// 		},
// 		{
// 			"ChangeImage",
// 			http.MethodPatch,
// 			"/api/v1/api/v1/images/:id/changeImage",
// 			handleFunctions.ImagesAPI.ChangeImage,
// 		},
// 		{
// 			"DeleteImageById",
// 			http.MethodDelete,
// 			"/api/v1/api/v1/images/:id",
// 			handleFunctions.ImagesAPI.DeleteImageById,
// 		},
// 		{
// 			"SearchImageById",
// 			http.MethodGet,
// 			"/api/v1/api/v1/images/:id",
// 			handleFunctions.ImagesAPI.SearchImageById,
// 		},
// 		{
// 			"SearchProductImage",
// 			http.MethodGet,
// 			"/api/v1/api/v1/images/products/:id",
// 			handleFunctions.ImagesAPI.SearchProductImage,
// 		},
// 		{
// 			"AddProduct",
// 			http.MethodPost,
// 			"/api/v1/api/v1/products",
// 			handleFunctions.ProductsAPI.AddProduct,
// 		},
// 		{
// 			"DecreaseParametr",
// 			http.MethodPatch,
// 			"/api/v1/api/v1/products/:id/decrease",
// 			handleFunctions.ProductsAPI.DecreaseParametr,
// 		},
// 		{
// 			"DeleteProductById",
// 			http.MethodDelete,
// 			"/api/v1/api/v1/products/:id",
// 			handleFunctions.ProductsAPI.DeleteProductById,
// 		},
// 		{
// 			"GetAllProduct",
// 			http.MethodGet,
// 			"/api/v1/api/v1/products",
// 			handleFunctions.ProductsAPI.GetAllProduct,
// 		},
// 		{
// 			"SearchProductById",
// 			http.MethodGet,
// 			"/api/v1/api/v1/products/:id",
// 			handleFunctions.ProductsAPI.SearchProductById,
// 		},
// 		{
// 			"AddSupplier",
// 			http.MethodPost,
// 			"/api/v1/api/v1/suppliers",
// 			handleFunctions.SuppliersAPI.AddSupplier,
// 		},
// 		{
// 			"ChangeAddressParametr",
// 			http.MethodPatch,
// 			"/api/v1/api/v1/suppliers/:id/changeAddress",
// 			handleFunctions.SuppliersAPI.ChangeAddressParametr,
// 		},
// 		{
// 			"DeleteSupplierById",
// 			http.MethodDelete,
// 			"/api/v1/api/v1/suppliers/:id",
// 			handleFunctions.SuppliersAPI.DeleteSupplierById,
// 		},
// 		{
// 			"GetAllSuppliers",
// 			http.MethodGet,
// 			"/api/v1/api/v1/suppliers",
// 			handleFunctions.SuppliersAPI.GetAllSuppliers,
// 		},
// 		{
// 			"SearchSupplierById",
// 			http.MethodGet,
// 			"/api/v1/api/v1/suppliers/:id",
// 			handleFunctions.SuppliersAPI.SearchSupplierById,
// 		},
// 	}
// }

package api

import (
	db "gilanggsb/simplebank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store db.Store
	route *gin.Engine
}

type BaseResponse[T any] struct {
	Status  int    `json:"status"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

var BasePath = "/api/v1"

// NewServer creates new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", ValidCurrency)
	}

	// add routes to router
	// Group with basePath
	api := router.Group(BasePath)
	{
		api.POST("/accounts", server.CreateAccount)
		api.GET("/accounts/:id", server.GetAccount)
		api.GET("/accounts", server.ListAccount)

		api.POST("/transfers", server.CreateTransfer)
	}

	// Handle 404
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, "route not found")
	})

	server.route = router
	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.route.Run(address)
}

func baseResponse[T any](baseResponse *BaseResponse[T]) gin.H {
	return gin.H{
		"status":  baseResponse.Status,
		"message": baseResponse.Message,
		"data":    baseResponse.Data,
	}
}

// Error wrapper
func errorResponse(err error, status int) gin.H {
	return gin.H{
		"status":  status,
		"message": err.Error(),
	}
}

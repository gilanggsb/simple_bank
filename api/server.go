package api

import (
	db "gilanggsb/simplebank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store *db.Store
	route *gin.Engine
}

type BaseResponse[T any] struct {
	Status  int    `json:"status"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

// NewServer creates new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	basePath := "/api/v1"
	server := &Server{store: store}
	router := gin.Default()

	// add routes to router
	// Group with basePath
	api := router.Group(basePath)
	{
		api.POST("/accounts", server.CreateAccount)
		api.GET("/accounts/:id", server.GetAccount)
		api.GET("/accounts", server.ListAccount)
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

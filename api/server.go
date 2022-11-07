package api

import (
	"fmt"

	db "github.com/bagashiz/Simple-Bank/db/sqlc"
	"github.com/bagashiz/Simple-Bank/token"
	"github.com/bagashiz/Simple-Bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for banking service.
type Server struct {
  config     util.Config
	store      db.Store
  tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
  tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
  if err != nil {
    return nil, fmt.Errorf("cannot create token master: %v", err)
  }
	server := &Server{
    config:     config,
    store:      store,
    tokenMaker: tokenMaker,
  }
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.POST("/transfers", server.createTransfer)
	router.POST("/users", server.createUser)

	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccounts)

	server.router = router

	return server, nil
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse is a common format for API errors.
func errorReponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

package server

import (
	"net/http"

	"github.com/anton-uvarenko/promova_test/internal/transport"
	"github.com/gin-gonic/gin"
)

func NewServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
}

func SetUpRoutes(handler *transport.Handler) http.Handler {
	router := gin.New()
	return router
}

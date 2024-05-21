package main

import (
	"github.com/anton-uvarenko/promova_test/internal/pkg/server"
	"github.com/anton-uvarenko/promova_test/internal/service"
	"github.com/anton-uvarenko/promova_test/internal/transport"
)

func main() {
	appService := service.NewService()
	handler := transport.NewHandler(appService)

	router := server.SetUpRoutes(handler)
	httpServer := server.NewServer(router)

	httpServer.ListenAndServe()
}

package main

import (
	"github.com/anton-uvarenko/promova_test/internal/core"
	"github.com/anton-uvarenko/promova_test/internal/db"
	"github.com/anton-uvarenko/promova_test/internal/pkg/server"
	"github.com/anton-uvarenko/promova_test/internal/service"
	"github.com/anton-uvarenko/promova_test/internal/transport"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	conn := db.Connect()
	repo := core.New(conn)

	appService := service.NewService(repo)
	handler := transport.NewHandler(appService.NewsService)

	router := server.SetUpRoutes(handler.NewsHandler)
	httpServer := server.NewServer(router, "8080")

	httpServer.ListenAndServe()
}

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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

	go httpServer.ListenAndServe()
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, syscall.SIGTERM)

	<-finish

	conn.Close(context.Background())
}

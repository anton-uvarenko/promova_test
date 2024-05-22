package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewServer(handler http.Handler, addr string) *http.Server {
	return &http.Server{
		Addr:    ":" + addr,
		Handler: handler,
	}
}

type newsHandler interface {
	AddNews(ctx *gin.Context)
	UpdateNews(ctx *gin.Context)
	GetNewsById(ctx *gin.Context)
	GetAllNews(ctx *gin.Context)
	DeleteNews(ctx *gin.Context)
}

func SetUpRoutes(newsHandler newsHandler) http.Handler {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	router.POST("/posts", newsHandler.AddNews)
	router.GET("/posts", newsHandler.GetAllNews)
	router.PUT("/posts/:id", newsHandler.UpdateNews)
	router.GET("/posts/:id", newsHandler.GetNewsById)
	router.DELETE("/posts/:id", newsHandler.DeleteNews)

	return router
}

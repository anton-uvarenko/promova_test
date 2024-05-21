package transport

import "github.com/anton-uvarenko/promova_test/internal/service"

type Handler struct{}

func NewHandler(service *service.Service) *Handler {
	return &Handler{}
}

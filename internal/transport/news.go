package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/anton-uvarenko/promova_test/internal/core"
	"github.com/anton-uvarenko/promova_test/internal/pkg"
	"github.com/anton-uvarenko/promova_test/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type NewsHandler struct {
	newsService newsService
}

func NewNewsHandler(newsService newsService) *NewsHandler {
	return &NewsHandler{
		newsService: newsService,
	}
}

type newsService interface {
	AddNews(ctx context.Context, params core.AddNewsParams) (int32, error)
}

func (h *NewsHandler) AddNews(ctx *gin.Context) {
	type AddNewsPayload struct {
		Title   string `json:"title" binding:"required,gt=2,lt=50"`
		Content string `json:"content" binding:"required"`
	}

	var payload AddNewsPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Code:  response.InvalidPayload,
			Error: fmt.Errorf("%w: [%w]", pkg.ErrInvalidPayload, err).Error(),
		})
		return
	}

	id, err := h.newsService.AddNews(ctx, core.AddNewsParams{
		Title:   pgtype.Text{String: payload.Title, Valid: true},
		Content: pgtype.Text{String: payload.Content, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pkg.ErrEntityAlreadyExists) {
			ctx.AbortWithStatusJSON(http.StatusConflict, response.Response{
				Code:  response.EntityAlreadyExists,
				Error: err.Error(),
			})
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: response.Ok,
		Data: response.AddNewsData{
			Id: int(id),
		},
	})
}

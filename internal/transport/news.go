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
	UpdatNews(ctx context.Context, params core.UpdateNewsParams) error
	GetNewsById(ctx context.Context, id int32) (core.News, error)
	GetAllNews(ctx context.Context) ([]core.News, error)
	DeleteNews(ctx context.Context, id int32) error
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

func (h *NewsHandler) UpdateNews(ctx *gin.Context) {
	type UpdateNewsPayload struct {
		Title   string `json:"title" binding:"required,gt=2,lt=50"`
		Content string `json:"content" binding:"required"`
	}

	var payload UpdateNewsPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Code:  response.InvalidPayload,
			Error: fmt.Errorf("%w: [%w]", pkg.ErrInvalidPayload, err).Error(),
		})
		return
	}

	type UriPayload struct {
		Id int `uri:"id"`
	}
	var uriPayload UriPayload
	err = ctx.ShouldBindUri(&uriPayload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Code:  response.InvalidPayload,
			Error: pkg.ErrInvalidUriParameters.Error(),
		})
		return
	}

	err = h.newsService.UpdatNews(ctx, core.UpdateNewsParams{
		ID:      int32(uriPayload.Id),
		Title:   pgtype.Text{String: payload.Title, Valid: true},
		Content: pgtype.Text{String: payload.Content, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, response.Response{
				Code:  response.NotFound,
				Error: pkg.ErrNotFound.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
			Code:  response.InternalError,
			Error: pkg.ErrDbInternal.Error(),
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *NewsHandler) GetNewsById(ctx *gin.Context) {
	type UriPayload struct {
		Id int `uri:"id"`
	}
	var uriPayload UriPayload
	err := ctx.ShouldBindUri(&uriPayload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Code:  response.InvalidPayload,
			Error: pkg.ErrInvalidUriParameters.Error(),
		})
		return
	}

	news, err := h.newsService.GetNewsById(ctx, int32(uriPayload.Id))
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, response.Response{
				Code:  response.NotFound,
				Error: pkg.ErrNotFound.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: response.Ok,
		Data: news,
	})
}

func (h *NewsHandler) GetAllNews(ctx *gin.Context) {
	news, err := h.newsService.GetAllNews(ctx)
	if err != nil {
		if errors.Is(err, pkg.ErrNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, response.Response{
				Code:  response.NotFound,
				Error: pkg.ErrNotFound.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: response.Ok,
		Data: news,
	})
}

func (h *NewsHandler) DeleteNews(ctx *gin.Context) {
	type UriPayload struct {
		Id int `uri:"id"`
	}
	var uriPayload UriPayload
	err := ctx.BindUri(&uriPayload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Response{
			Code:  response.InvalidPayload,
			Error: pkg.ErrInvalidPayload.Error(),
		})
		return
	}

	err = h.newsService.DeleteNews(ctx, int32(uriPayload.Id))
	if err != nil {
		if errors.Is(err, pkg.ErrEntityAlreadyDeleted) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, response.Response{
				Code:  response.NotFound,
				Error: pkg.ErrEntityAlreadyDeleted.Error(),
			})
			return
		}
	}

	ctx.Status(http.StatusOK)
}

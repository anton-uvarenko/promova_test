package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/anton-uvarenko/promova_test/internal/core"
	"github.com/anton-uvarenko/promova_test/internal/pkg"
	"github.com/jackc/pgx/v5/pgconn"
)

type NewsService struct {
	newsRepo newsRepo
}

func NewNewsService(newsRepo newsRepo) *NewsService {
	return &NewsService{
		newsRepo: newsRepo,
	}
}

type newsRepo interface {
	AddNews(ctx context.Context, arg core.AddNewsParams) (int32, error)
	DeleteNews(ctx context.Context, id int32) error
	GetAllNews(ctx context.Context) ([]core.News, error)
	GetNewsById(ctx context.Context, id int32) (core.News, error)
	UpdateNews(ctx context.Context, arg core.UpdateNewsParams) error
}

func (s *NewsService) AddNews(ctx context.Context, params core.AddNewsParams) (int32, error) {
	id, err := s.newsRepo.AddNews(ctx, params)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			err := err.(*pgconn.PgError)
			// duplicate key error
			if err.Code == "23505" {
				return 0, pkg.ErrEntityAlreadyExists
			}
		}

		fmt.Printf("%v: [%v]", pkg.ErrDbInternal, err)
		return 0, fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	return id, nil
}

package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/anton-uvarenko/promova_test/internal/core"
	"github.com/anton-uvarenko/promova_test/internal/pkg"
	"github.com/jackc/pgx/v5"
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

		fmt.Printf("%v: [%v]\n", pkg.ErrDbInternal, err)
		return 0, fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	return id, nil
}

func (s *NewsService) UpdatNews(ctx context.Context, params core.UpdateNewsParams) error {
	err := s.newsRepo.UpdateNews(ctx, params)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			err := err.(*pgconn.PgError)
			if err.Code == "23505" {
				return pkg.ErrEntityAlreadyExists
			}
		}

		fmt.Printf("%v: [%v]\n", pkg.ErrDbInternal, err)
		return fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}
	return nil
}

func (s *NewsService) GetAllNews(ctx context.Context) ([]core.News, error) {
	news, err := s.newsRepo.GetAllNews(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}
		fmt.Printf("%v: [%v]\n", pkg.ErrDbInternal, err)
		return nil, fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	return news, nil
}

func (s *NewsService) GetNewsById(ctx context.Context, id int32) (core.News, error) {
	news, err := s.newsRepo.GetNewsById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.News{}, pkg.ErrNotFound
		}

		fmt.Printf("%v: [%v]\n", pkg.ErrDbInternal, err)
		return core.News{}, fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	return news, nil
}

func (s *NewsService) DeleteNews(ctx context.Context, id int32) error {
	_, err := s.GetNewsById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pkg.ErrEntityAlreadyDeleted
		}

		fmt.Printf("%v: [%v]\n", pkg.ErrDbInternal, err)
		return pkg.ErrDbInternal
	}

	err = s.newsRepo.DeleteNews(ctx, id)
	if err != nil {
		fmt.Printf("%v: [%v]\n", pkg.ErrDbInternal, err)
		return fmt.Errorf("%w: [%w]", pkg.ErrDbInternal, err)
	}

	return nil
}

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/anton-uvarenko/promova_test/internal/core"
	"github.com/anton-uvarenko/promova_test/internal/pkg"
	"github.com/go-playground/assert/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type NewsRepoMock struct {
	ErrAddNewsToReturn     error
	ErrUpdateNewsToReturn  error
	ErrGetAllNewsToReturn  error
	ErrGetNewsByIdToReturn error
	ErrDeleteNewsToReturn  error
}

func (m *NewsRepoMock) AddNews(ctx context.Context, arg core.AddNewsParams) (int32, error) {
	if m.ErrAddNewsToReturn != nil {
		return 0, m.ErrAddNewsToReturn
	}
	return 1, nil
}

func (m *NewsRepoMock) UpdateNews(ctx context.Context, arg core.UpdateNewsParams) error {
	if m.ErrUpdateNewsToReturn != nil {
		return m.ErrUpdateNewsToReturn
	}
	return nil
}

func (m *NewsRepoMock) GetAllNews(ctx context.Context) ([]core.News, error) {
	if m.ErrGetAllNewsToReturn != nil {
		return nil, m.ErrGetAllNewsToReturn
	}
	return []core.News{
		{
			ID: 1,
		},
	}, nil
}

func (m *NewsRepoMock) GetNewsById(ctx context.Context, id int32) (core.News, error) {
	if m.ErrGetNewsByIdToReturn != nil {
		return core.News{}, m.ErrGetNewsByIdToReturn
	}
	return core.News{
		ID: 1,
	}, nil
}

func (m *NewsRepoMock) DeleteNews(ctx context.Context, id int32) error {
	if m.ErrDeleteNewsToReturn != nil {
		return m.ErrDeleteNewsToReturn
	}
	return nil
}

func TestAddNews(t *testing.T) {
	repo := &NewsRepoMock{}
	service := NewNewsService(repo)

	testTable := []struct {
		Name                string
		ErrRepoShouldReturn error
		ExpectedError       error
		ExpectedResult      int32
	}{
		{
			Name:           "Ok",
			ExpectedError:  nil,
			ExpectedResult: 1,
		},
		{
			Name: "Err entity already exist",
			ErrRepoShouldReturn: &pgconn.PgError{
				Code: "23505",
			},
			ExpectedError:  pkg.ErrEntityAlreadyExists,
			ExpectedResult: 0,
		},
		{
			Name:                "Err db internal",
			ErrRepoShouldReturn: errors.New("some unexpected error"),
			ExpectedError:       pkg.ErrDbInternal,
			ExpectedResult:      0,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			repo.ErrAddNewsToReturn = testCase.ErrRepoShouldReturn

			id, err := service.AddNews(context.Background(), core.AddNewsParams{})

			assert.Equal(t, errors.Is(err, testCase.ExpectedError), true)
			assert.Equal(t, id, testCase.ExpectedResult)
		})
	}
}

func TestUpdatNews(t *testing.T) {
	repo := &NewsRepoMock{}
	service := NewNewsService(repo)

	testTable := []struct {
		Name                string
		ErrRepoShouldReturn error
		ExpectedError       error
	}{
		{
			Name:                "Ok",
			ErrRepoShouldReturn: nil,
			ExpectedError:       nil,
		},
		{
			Name: "Err duplicate key",
			ErrRepoShouldReturn: &pgconn.PgError{
				Code: "23505",
			},
			ExpectedError: pkg.ErrEntityAlreadyExists,
		},
		{
			Name:                "Err db internal",
			ErrRepoShouldReturn: errors.New("some unexpected error"),
			ExpectedError:       pkg.ErrDbInternal,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			repo.ErrUpdateNewsToReturn = testCase.ErrRepoShouldReturn
			err := service.UpdatNews(context.Background(), core.UpdateNewsParams{})
			assert.Equal(t, errors.Is(err, testCase.ExpectedError), true)
		})
	}
}

func TestGetAllNews(t *testing.T) {
	repo := &NewsRepoMock{}
	service := NewNewsService(repo)

	testTable := []struct {
		Name                string
		ErrRepoShouldReturn error
		ExpectedError       error
		ExpectedResult      []core.News
	}{
		{
			Name:                "Ok",
			ErrRepoShouldReturn: nil,
			ExpectedError:       nil,
			ExpectedResult: []core.News{
				{
					ID: 1,
				},
			},
		},
		{
			Name:                "Err not found",
			ErrRepoShouldReturn: pgx.ErrNoRows,
			ExpectedError:       pkg.ErrNotFound,
			ExpectedResult:      nil,
		},
		{
			Name:                "Err db internal",
			ErrRepoShouldReturn: errors.New("some unexpected error"),
			ExpectedError:       pkg.ErrDbInternal,
			ExpectedResult:      nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			repo.ErrGetAllNewsToReturn = testCase.ErrRepoShouldReturn
			result, err := service.GetAllNews(context.Background())

			assert.Equal(t, errors.Is(err, testCase.ExpectedError), true)
			if len(testCase.ExpectedResult) > 0 {
				assert.Equal(t, result[0].ID, testCase.ExpectedResult[0].ID)
			}
		})
	}
}

func TestGetNewsById(t *testing.T) {
	repo := &NewsRepoMock{}
	service := NewNewsService(repo)
	testTable := []struct {
		Name                string
		ErrRepoShouldReturn error
		ExpectedError       error
		ExpectedResult      core.News
	}{
		{
			Name:                "Ok",
			ErrRepoShouldReturn: nil,
			ExpectedError:       nil,
			ExpectedResult: core.News{
				ID: 1,
			},
		},
		{
			Name:                "Err not found",
			ErrRepoShouldReturn: pgx.ErrNoRows,
			ExpectedError:       pkg.ErrNotFound,
			ExpectedResult:      core.News{},
		},
		{
			Name:                "Err db internal",
			ErrRepoShouldReturn: errors.New("some unexpected error"),
			ExpectedError:       pkg.ErrDbInternal,
			ExpectedResult:      core.News{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			repo.ErrGetNewsByIdToReturn = testCase.ErrRepoShouldReturn

			result, err := service.GetNewsById(context.Background(), 1)

			assert.Equal(t, errors.Is(err, testCase.ExpectedError), true)
			assert.Equal(t, result.ID, testCase.ExpectedResult.ID)
		})
	}
}

func TestDeleteNews(t *testing.T) {
	repo := &NewsRepoMock{}
	service := NewNewsService(repo)
	testTable := []struct {
		Name                       string
		ErrDelteNewsShouldReturn   error
		ErrGetNewsByIdShouldReturn error
		ExpectedError              error
	}{
		{
			Name:                       "Ok",
			ErrDelteNewsShouldReturn:   nil,
			ErrGetNewsByIdShouldReturn: nil,
			ExpectedError:              nil,
		},
		{
			Name:                       "Err already deleted",
			ErrDelteNewsShouldReturn:   nil,
			ErrGetNewsByIdShouldReturn: pgx.ErrNoRows,
			ExpectedError:              pkg.ErrEntityAlreadyDeleted,
		},
		{
			Name:                       "Err db internal",
			ErrDelteNewsShouldReturn:   errors.New(`some unexpected err`),
			ErrGetNewsByIdShouldReturn: nil,
			ExpectedError:              pkg.ErrDbInternal,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			repo.ErrDeleteNewsToReturn = testCase.ErrDelteNewsShouldReturn
			repo.ErrGetNewsByIdToReturn = testCase.ErrGetNewsByIdShouldReturn

			err := service.DeleteNews(context.Background(), 1)
			assert.Equal(t, errors.Is(err, testCase.ExpectedError), true)
		})
	}
}

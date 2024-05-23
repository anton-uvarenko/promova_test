package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/anton-uvarenko/promova_test/internal/core"
	"github.com/anton-uvarenko/promova_test/internal/pkg"
	"github.com/anton-uvarenko/promova_test/internal/pkg/payload"
	"github.com/anton-uvarenko/promova_test/internal/pkg/response"
	"github.com/anton-uvarenko/promova_test/internal/pkg/server"
	"github.com/go-playground/assert/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type newsServiceMock struct {
	ErrAddNewsToReturn     error
	ErrUpdateNewsToReturn  error
	ErrGetNewsByIdToReturn error
	ErrGetAllNewsToReturn  error
	ErrDeleteNewsToReturn  error
}

func (m *newsServiceMock) AddNews(ctx context.Context, params core.AddNewsParams) (int32, error) {
	if m.ErrAddNewsToReturn != nil {
		return 0, m.ErrAddNewsToReturn
	}
	return 1, nil
}

func (m *newsServiceMock) UpdatNews(ctx context.Context, params core.UpdateNewsParams) error {
	if m.ErrAddNewsToReturn != nil {
		return m.ErrUpdateNewsToReturn
	}

	return nil
}

func (m *newsServiceMock) GetNewsById(ctx context.Context, id int32) (core.News, error) {
	if m.ErrGetNewsByIdToReturn != nil {
		return core.News{}, m.ErrGetNewsByIdToReturn
	}

	return core.News{
		ID:        id,
		Title:     pgtype.Text{String: "some title", Valid: true},
		Content:   pgtype.Text{String: "some content", Valid: true},
		CreatedAt: pgtype.Timestamp{},
	}, nil
}

func (m *newsServiceMock) GetAllNews(ctx context.Context) ([]core.News, error) {
	if m.ErrGetAllNewsToReturn != nil {
		return nil, m.ErrGetAllNewsToReturn
	}

	return []core.News{
		{
			ID:      1,
			Title:   pgtype.Text{String: "some title", Valid: true},
			Content: pgtype.Text{String: "some content", Valid: true},
		},
	}, nil
}

func (m *newsServiceMock) DeleteNews(ctx context.Context, id int32) error {
	if m.ErrDeleteNewsToReturn != nil {
		return m.ErrDeleteNewsToReturn
	}

	return nil
}

type AddNewsResponse struct {
	Code int                  `json:"code"`
	Data response.AddNewsData `json:"data"`
}

var (
	httpServer          http.Server
	newsServiceInstance *newsServiceMock
)

func TestMain(m *testing.M) {
	newsServiceInstance = &newsServiceMock{}
	handler := NewHandler(newsServiceInstance)
	router := server.SetUpRoutes(handler.NewsHandler)
	httpServer := server.NewServer(router, "8081")
	go httpServer.ListenAndServe()

	m.Run()
}

func TestAddNews(t *testing.T) {
	testTable := []struct {
		Name                     string
		RequestPayload           any
		ErrorServiceShouldReturn error
		ExpectedResult           AddNewsResponse
		ExpectedStatusCode       int
	}{
		{
			Name: "Ok",
			RequestPayload: payload.AddNewsPayload{
				Title:   "some title",
				Content: "some content",
			},
			ErrorServiceShouldReturn: nil,
			ExpectedResult: AddNewsResponse{
				Code: response.Ok,
				Data: response.AddNewsData{
					Id: 1,
				},
			},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:                     "Decode error",
			RequestPayload:           "some string that fail decoding",
			ErrorServiceShouldReturn: nil,
			ExpectedResult: AddNewsResponse{
				Code: response.InvalidPayload,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name: "Validation error",
			RequestPayload: payload.AddNewsPayload{
				Title:   "f",
				Content: "some content",
			},
			ErrorServiceShouldReturn: nil,
			ExpectedResult: AddNewsResponse{
				Code: response.InvalidPayload,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name: "Entity already exist",
			RequestPayload: payload.AddNewsPayload{
				Title:   "some title",
				Content: "some content",
			},
			ErrorServiceShouldReturn: pkg.ErrEntityAlreadyExists,
			ExpectedResult: AddNewsResponse{
				Code: response.EntityAlreadyExists,
			},
			ExpectedStatusCode: http.StatusConflict,
		},
		{
			Name: "Db Internal error",
			RequestPayload: payload.AddNewsPayload{
				Title:   "some title",
				Content: "some content",
			},
			ExpectedResult: AddNewsResponse{
				Code: response.InternalError,
			},
			ErrorServiceShouldReturn: pkg.ErrDbInternal,
			ExpectedStatusCode:       http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			newsServiceInstance.ErrAddNewsToReturn = testCase.ErrorServiceShouldReturn
			pl, _ := json.Marshal(testCase.RequestPayload)

			r, _ := http.NewRequest(http.MethodPost, "http://localhost:8081/posts", bytes.NewBuffer(pl))
			resp, _ := http.DefaultClient.Do(r)

			assert.Equal(t, resp.StatusCode, testCase.ExpectedStatusCode)

			body, _ := io.ReadAll(resp.Body)
			respResult := AddNewsResponse{}
			err := json.Unmarshal(body, &respResult)
			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}

			assert.Equal(t, respResult.Code, testCase.ExpectedResult.Code)

			if testCase.ExpectedStatusCode == http.StatusOK {
				assert.Equal(t, respResult.Data.Id, testCase.ExpectedResult.Data.Id)
			}
		})
	}
}

func TestUpdateNews(t *testing.T) {
	testTable := []struct {
		Name                     string
		RequestPayload           any
		UriParam                 string
		ErrorServiceShouldReturn error
		ExpectedResult           response.Response
		ExpectedStatusCode       int
	}{
		{
			Name: "Ok",
			RequestPayload: payload.UpdateNewsPayload{
				Title:   "some title",
				Content: "some content",
			},
			UriParam:                 "1",
			ErrorServiceShouldReturn: nil,
			ExpectedResult: response.Response{
				Code: response.Ok,
			},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:                     "Error decode payload",
			RequestPayload:           "adsfadsf",
			UriParam:                 "1",
			ErrorServiceShouldReturn: nil,
			ExpectedResult: response.Response{
				Code: response.InvalidPayload,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name: "Error validation",
			RequestPayload: payload.UpdateNewsPayload{
				Title:   "s",
				Content: "some content",
			},
			UriParam:                 "1",
			ErrorServiceShouldReturn: nil,
			ExpectedResult: response.Response{
				Code: response.InvalidPayload,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name: "Invalid uri param",
			RequestPayload: payload.UpdateNewsPayload{
				Title:   "some title",
				Content: "some content",
			},
			UriParam:                 "aa",
			ErrorServiceShouldReturn: nil,
			ExpectedResult: response.Response{
				Code: response.InvalidPayload,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name: "Erorr not found",
			RequestPayload: payload.UpdateNewsPayload{
				Title:   "some title",
				Content: "some content",
			},
			UriParam:                 "1",
			ErrorServiceShouldReturn: pkg.ErrNotFound,
			ExpectedResult: response.Response{
				Code: response.NotFound,
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name: "Error entity already exists",
			RequestPayload: payload.UpdateNewsPayload{
				Title:   "some title",
				Content: "some content",
			},
			UriParam:                 "1",
			ErrorServiceShouldReturn: pkg.ErrEntityAlreadyExists,
			ExpectedResult: response.Response{
				Code: response.EntityAlreadyExists,
			},
			ExpectedStatusCode: http.StatusConflict,
		},
		{
			Name: "Error db internal",
			RequestPayload: payload.UpdateNewsPayload{
				Title:   "some title",
				Content: "some content",
			},
			UriParam:                 "1",
			ErrorServiceShouldReturn: pkg.ErrDbInternal,
			ExpectedResult: response.Response{
				Code: response.InternalError,
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			newsServiceInstance.ErrUpdateNewsToReturn = testCase.ErrorServiceShouldReturn
			pl, _ := json.Marshal(testCase.RequestPayload)

			r, _ := http.NewRequest(http.MethodPut, "http://localhost:8081/posts/"+testCase.UriParam, bytes.NewBuffer(pl))
			resp, _ := http.DefaultClient.Do(r)

			assert.Equal(t, resp.StatusCode, testCase.ExpectedStatusCode)

			var respResult response.Response
			err := json.NewDecoder(resp.Body).Decode(&respResult)
			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}

			assert.Equal(t, respResult.Code, testCase.ExpectedResult.Code)
		})
	}
}

type GetNewsByIdResponse struct {
	Code int               `json:"code"`
	Data response.NewsData `json:"data"`
}

func TestGetNewsById(t *testing.T) {
	testTable := []struct {
		Name                     string
		UriParam                 string
		ErrorServiceShouldReturn error
		ExpectedResult           GetNewsByIdResponse
		ExpectedStatusCode       int
	}{
		{
			Name:                     "Ok",
			UriParam:                 "1",
			ErrorServiceShouldReturn: nil,
			ExpectedResult: GetNewsByIdResponse{
				Code: response.Ok,
				Data: response.NewsData{
					Id:      1,
					Title:   "some title",
					Content: "some content",
				},
			},
			ExpectedStatusCode: 200,
		},
		{
			Name:                     "Invalid uri param",
			UriParam:                 "adsfasd",
			ErrorServiceShouldReturn: pkg.ErrInvalidPayload,
			ExpectedResult: GetNewsByIdResponse{
				Code: response.InvalidPayload,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:                     "Error not found",
			UriParam:                 "1",
			ErrorServiceShouldReturn: pkg.ErrNotFound,
			ExpectedResult: GetNewsByIdResponse{
				Code: response.NotFound,
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:                     "Error db internal",
			UriParam:                 "1",
			ErrorServiceShouldReturn: pkg.ErrDbInternal,
			ExpectedResult: GetNewsByIdResponse{
				Code: response.InternalError,
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			newsServiceInstance.ErrGetNewsByIdToReturn = testCase.ErrorServiceShouldReturn

			r, _ := http.NewRequest(http.MethodGet, "http://localhost:8081/posts/"+testCase.UriParam, nil)
			resp, _ := http.DefaultClient.Do(r)

			assert.Equal(t, resp.StatusCode, testCase.ExpectedStatusCode)

			var respResult GetNewsByIdResponse
			err := json.NewDecoder(resp.Body).Decode(&respResult)
			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}

			assert.Equal(t, respResult.Code, testCase.ExpectedResult.Code)
			assert.Equal(t, respResult.Data.Id, testCase.ExpectedResult.Data.Id)
			assert.Equal(t, respResult.Data.Title, testCase.ExpectedResult.Data.Title)
			assert.Equal(t, respResult.Data.Content, testCase.ExpectedResult.Data.Content)
		})
	}
}

type GetAllNewsResponse struct {
	Code int                 `json:"code"`
	Data []response.NewsData `json:"data"`
}

func TestGetAllNews(t *testing.T) {
	testTable := []struct {
		Name                     string
		ErrorServiceShouldReturn error
		ExpectedResult           GetAllNewsResponse
		ExpectedStatusCode       int
	}{
		{
			Name:                     "Ok",
			ErrorServiceShouldReturn: nil,
			ExpectedResult: GetAllNewsResponse{
				Code: response.Ok,
				Data: []response.NewsData{
					{
						Id:      1,
						Title:   "some title",
						Content: "some content",
					},
				},
			},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:                     "Error not found",
			ErrorServiceShouldReturn: pkg.ErrNotFound,
			ExpectedResult: GetAllNewsResponse{
				Code: response.NotFound,
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Name:                     "Error db internal",
			ErrorServiceShouldReturn: pkg.ErrDbInternal,
			ExpectedResult: GetAllNewsResponse{
				Code: response.InternalError,
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			newsServiceInstance.ErrGetAllNewsToReturn = testCase.ErrorServiceShouldReturn

			r, _ := http.NewRequest(http.MethodGet, "http://localhost:8081/posts", nil)
			resp, _ := http.DefaultClient.Do(r)

			assert.Equal(t, resp.StatusCode, testCase.ExpectedStatusCode)

			var respResult GetAllNewsResponse
			err := json.NewDecoder(resp.Body).Decode(&respResult)
			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}

			if resp.StatusCode == http.StatusOK {
				assert.Equal(t, respResult.Code, testCase.ExpectedResult.Code)
				assert.Equal(t, respResult.Data[0].Id, testCase.ExpectedResult.Data[0].Id)
				assert.Equal(t, respResult.Data[0].Title, testCase.ExpectedResult.Data[0].Title)
				assert.Equal(t, respResult.Data[0].Content, testCase.ExpectedResult.Data[0].Content)
			}
		})
	}
}

func TestDeleteNews(t *testing.T) {
	testTable := []struct {
		Name                     string
		UriParam                 string
		ErrorServiceShouldReturn error
		ExpectedResult           response.Response
		ExpectedStatusCode       int
	}{
		{
			Name:                     "Ok",
			UriParam:                 "1",
			ErrorServiceShouldReturn: nil,
			ExpectedResult: response.Response{
				Code: response.Ok,
			},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:                     "Err invalied payload",
			UriParam:                 "adfasdf",
			ErrorServiceShouldReturn: nil,
			ExpectedResult: response.Response{
				Code: response.InvalidPayload,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:                     "Err entity already deleted",
			UriParam:                 "1",
			ErrorServiceShouldReturn: pkg.ErrEntityAlreadyDeleted,
			ExpectedResult:           response.Response{Code: response.NotFound},
			ExpectedStatusCode:       http.StatusNotFound,
		},
		{
			Name:                     "Err internal",
			UriParam:                 "1",
			ErrorServiceShouldReturn: pkg.ErrDbInternal,
			ExpectedResult: response.Response{
				Code: response.InternalError,
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			newsServiceInstance.ErrDeleteNewsToReturn = testCase.ErrorServiceShouldReturn

			r, _ := http.NewRequest(http.MethodDelete, "http://localhost:8081/posts/"+testCase.UriParam, nil)
			resp, _ := http.DefaultClient.Do(r)

			assert.Equal(t, resp.StatusCode, testCase.ExpectedStatusCode)

			var respResult GetAllNewsResponse
			err := json.NewDecoder(resp.Body).Decode(&respResult)
			if err != nil {
				t.Error(err)
				t.Fail()
				return
			}

			assert.Equal(t, respResult.Code, testCase.ExpectedResult.Code)
		})
	}
}

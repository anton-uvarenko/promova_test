package service

type Service struct {
	NewsService *NewsService
}

func NewService(newsRepo newsRepo) *Service {
	return &Service{
		NewsService: NewNewsService(newsRepo),
	}
}

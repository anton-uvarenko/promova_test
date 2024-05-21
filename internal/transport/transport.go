package transport

type Handler struct {
	NewsHandler *NewsHandler
}

func NewHandler(newsService newsService) *Handler {
	return &Handler{
		NewsHandler: NewNewsHandler(newsService),
	}
}

package response

type AddNewsData struct {
	Id int `json:"id"`
}

type NewsData struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

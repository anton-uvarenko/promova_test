package payload

type AddNewsPayload struct {
	Title   string `json:"title" binding:"required,gt=2,lt=50"`
	Content string `json:"content" binding:"required"`
}

type UpdateNewsPayload struct {
	Title   string `json:"title" binding:"required,gt=2,lt=50"`
	Content string `json:"content" binding:"required"`
}

type IdUriPayload struct {
	Id int `uri:"id"`
}

package dto

type CreateTodoPayload struct {
	Title   string `json:"title" validate:"required,max=666"`
	Content string `json:"content" validate:"required,max=6666"`
}

type UpdateTodoPayload struct {
	CreateTodoPayload
	Done bool `json:"done" validate:"required"`
}

package model

type CreateTodoRequest struct {
	Title  string `json:"title" binding:"required"`
	UserID int64
}

func (r CreateTodoRequest) ToTodo() Todo {
	return Todo{
		Title:  r.Title,
		UserID: r.UserID,
	}
}

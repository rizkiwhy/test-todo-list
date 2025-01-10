package presenter

type TodoRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

const (
	CreateTodoUnauthorizedErrorMessage   = "You are not authorized to create a todo"
	CreateTodoInvalidRequestMessage      = "Invalid create todo request"
	CreateTodoInternalServerErrorMessage = "Internal server error"
	CreateTodoSuccessMessage             = "Todo created successfully"
)

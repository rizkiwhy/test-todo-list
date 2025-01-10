package presenter

import (
	"net/http"
	"rizkiwhy/test-todo-list/package/user/model"
)

const (
	RegisterSuccessMessage         = "User registered successfully"
	RegisterFailureMessage         = "Failed to register user"
	RegisterInvalidRequestMessage  = "Invalid registration request"
	LoginSuccessMessage            = "User logged in successfully"
	LoginFailureMessage            = "Failed to log in user"
	LoginInvalidCredentialsMessage = "Invalid credentials"
)

var RegisterStatusCodeMap = map[string]int{
	model.ErrEmailAlreadyExists: http.StatusConflict,
	model.ErrNotFound:           http.StatusNotFound,
	model.ErrInternalError:      http.StatusInternalServerError,
}

var LoginStatusCodeMap = map[string]int{
	model.ErrNotFound:             http.StatusNotFound,
	model.ErrUnauthorizedAccess:   http.StatusUnauthorized,
	model.ErrInvalidEmailPassword: http.StatusUnauthorized,
	model.ErrInternalError:        http.StatusInternalServerError,
}

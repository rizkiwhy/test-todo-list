package presenter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(message string, data interface{}) gin.H {
	response := gin.H{
		"message": message,
		"status":  "success",
	}
	if data != nil {
		response["data"] = data
	}
	return response
}

func FailureResponse(title, message string) gin.H {
	return gin.H{
		"title":   title,
		"message": message,
		"status":  "error",
	}
}

func HandleError(c *gin.Context, err error, statusCodeMap map[string]int, title string) {
	httpStatusCode := http.StatusInternalServerError
	if code, ok := statusCodeMap[err.Error()]; ok {
		httpStatusCode = code
	}

	c.JSON(httpStatusCode, FailureResponse(title, err.Error()))
}

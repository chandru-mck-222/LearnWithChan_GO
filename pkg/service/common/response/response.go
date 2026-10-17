package response

import (
	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Status  string `json:"status"`
	Message string `json:"stirng"`
	Data    T      `json:"datat"`
}

func SendResponse[T any](ctx *gin.Context, status int, message, stat string, data T) {
	response := Response[T]{
		Status:  stat,
		Message: message,
		Data:    data,
	}
	ctx.JSON(status, response)
}

package common

import (
	"net/http"

	"chanlearning/pkg/service/common/response"

	"github.com/gin-gonic/gin"
)

// 400
func SendBadRequest(ctx *gin.Context, message string) {
	response.SendResponse[any](ctx, http.StatusBadRequest, message, "bad reuqest", nil)
}

// 200
func SendOkResponse(ctx *gin.Context, message string, Data any) {
	response.SendResponse(ctx, http.StatusOK, message, "success", Data)
}

// 201
func SendCreatedResponse(ctx *gin.Context, Message string, Data any) {
	response.SendResponse(ctx, http.StatusCreated, Message, "Created Successfully", Data)
}

// 401
func SendUnauthorizedResponse(ctx *gin.Context) {
	response.SendResponse[any](ctx, http.StatusUnauthorized, "unauthorized user", "", nil)

}

// 403
func SendForbidden(ctx *gin.Context) {
	response.SendResponse[any](ctx, http.StatusForbidden, "not followed message", "", nil)
}

// 409
func SendConflictResponse(ctx *gin.Context) {
	response.SendResponse[any](ctx, http.StatusConflict, "", "", nil)

}

// 413
func SendToomanyRequest(ctx *gin.Context) {
	response.SendResponse[any](ctx, http.StatusRequestEntityTooLarge, "too large file", "", nil)
}

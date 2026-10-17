package controllers

import (
	"chanlearning/pkg/domain/models"
	"chanlearning/pkg/service/common"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct{}

func (Ctrl UserController) CreateUser(ctx *gin.Context) {
	var body models.UserDetails
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		common.SendBadRequest(ctx, "not validate json")
	}
	if err := body.Validate(); err != nil {
		common.SendBadRequest(ctx, err.Error())
	}

}

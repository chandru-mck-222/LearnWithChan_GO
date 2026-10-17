package repository

import (
	"chanlearning/pkg/domain/models"

	"github.com/gin-gonic/gin"
)

type UserRepository struct {
}

func (Repo UserRepository) Create(ctx *gin.Context, body models.User) {

}

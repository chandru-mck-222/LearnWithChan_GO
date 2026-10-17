package repository

import "github.com/gin-gonic/gin"

type UserRepos interface {
	CreateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
}

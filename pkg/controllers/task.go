package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"chanlearning/pkg/domain/models"

	"github.com/gin-gonic/gin"
)

type TaskPublisher interface {
	Publish(context.Context, any) error
}

type TaskController struct {
	Publisher TaskPublisher
}

type createTaskRequest struct {
	Text string `json:"text" binding:"required"`
}

func (controller TaskController) Create(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 1<<20)
	var request createTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
		return
	}

	request.Text = strings.TrimSpace(request.Text)
	if request.Text == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
		return
	}
	if controller.Publisher == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "task queue unavailable"})
		return
	}

	id, err := newTaskID()
	if err != nil {
		log.Printf("create task ID: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create task"})
		return
	}
	task := models.TextTask{
		ID:        id,
		Text:      request.Text,
		CreatedAt: time.Now().UTC(),
	}
	if err := controller.Publisher.Publish(ctx.Request.Context(), task); err != nil {
		log.Printf("publish task %s: %v", task.ID, err)
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "could not enqueue task"})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"id":     task.ID,
		"status": "queued",
	})
}

func newTaskID() (string, error) {
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return "", fmt.Errorf("generate task ID: %w", err)
	}
	id[6] = (id[6] & 0x0f) | 0x40
	id[8] = (id[8] & 0x3f) | 0x80

	encoded := hex.EncodeToString(id[:])
	return encoded[:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:], nil
}

package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"chanlearning/pkg/domain/models"

	"github.com/gin-gonic/gin"
)

type taskPublisherFunc func(context.Context, any) error

func (publish taskPublisherFunc) Publish(ctx context.Context, payload any) error {
	return publish(ctx, payload)
}

func TestCreateTaskEnqueuesTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var published models.TextTask
	router := gin.New()
	router.POST("/tasks", (TaskController{
		Publisher: taskPublisherFunc(func(_ context.Context, payload any) error {
			task, ok := payload.(models.TextTask)
			if !ok {
				t.Fatalf("published payload has type %T, want models.TextTask", payload)
			}
			published = task
			return nil
		}),
	}).Create)

	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"text":"hello rabbitmq"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d; body: %s", recorder.Code, http.StatusAccepted, recorder.Body)
	}
	if published.ID == "" || published.Text != "hello rabbitmq" || published.CreatedAt.IsZero() {
		t.Fatalf("unexpected published task: %+v", published)
	}

	var response struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.ID != published.ID || response.Status != "queued" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestCreateTaskRejectsEmptyText(t *testing.T) {
	gin.SetMode(gin.TestMode)

	published := false
	router := gin.New()
	router.POST("/tasks", (TaskController{
		Publisher: taskPublisherFunc(func(context.Context, any) error {
			published = true
			return nil
		}),
	}).Create)

	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"text":"  "}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	if published {
		t.Fatal("empty task was published")
	}
}

func TestCreateTaskReturnsServiceUnavailableWhenPublishFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/tasks", (TaskController{
		Publisher: taskPublisherFunc(func(context.Context, any) error {
			return errors.New("broker unavailable")
		}),
	}).Create)

	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"text":"hello"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}
}

package server

import (
	"fmt"
	"os"
	"strings"

	"chanlearning/pkg/controllers"
	"chanlearning/pkg/messagebrokers/rabbitmq"
	"chanlearning/pkg/middleware"
	"chanlearning/pkg/service/common/envvar"

	"github.com/gin-gonic/gin"
)

func SetupServer() error {
	port := os.Getenv(envvar.PORT)
	if port == "" {
		return fmt.Errorf("%s environment variable is required", envvar.PORT)
	}

	conn, err := rabbitmq.RabbitDial()
	if err != nil {
		return err
	}
	defer conn.Close()

	queueName := os.Getenv(envvar.RabbitMQQueue)
	if queueName == "" {
		queueName = "text-tasks"
	}
	broker, err := conn.DeclareQueue(queueName, 0)
	if err != nil {
		return err
	}
	defer broker.Close()

	r := gin.Default()
	MiddlewareMapper(r)
	ServerRoutes(r, broker)
	return r.Run(":" + port)
}

func MiddlewareMapper(r *gin.Engine) {
	r.Use(middleware.Analystic())
	r.Use(gin.Recovery())
}

func ServerRoutes(r *gin.Engine, publisher controllers.TaskPublisher) {
	apiVersion := strings.Trim(os.Getenv(envvar.APIV), "/")
	if apiVersion == "" {
		apiVersion = "api/v1"
	}
	api := r.Group("/" + apiVersion)
	api.POST("/tasks", (controllers.TaskController{Publisher: publisher}).Create)
}

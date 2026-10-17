package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"chanlearning/pkg/messagebrokers/rabbitmq"
	"chanlearning/pkg/service/common/envvar"
	"chanlearning/pkg/tasks"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	client, err := rabbitmq.RabbitDial()
	if err != nil {
		return err
	}
	defer client.Close()

	queueName := os.Getenv(envvar.RabbitMQQueue)
	if queueName == "" {
		queueName = "text-tasks"
	}
	broker, err := client.DeclareQueue(queueName, 0)
	if err != nil {
		return err
	}

	log.Printf("worker listening on RabbitMQ queue %q", queueName)
	if err := broker.Consume(ctx, tasks.ProcessText); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

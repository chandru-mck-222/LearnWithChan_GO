package rabbitmq

import (
	"chanlearning/pkg/service/common/envvar"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn *amqp.Connection
}

type MessageBroker struct {
	queue   string
	channel *amqp.Channel
	mu      sync.Mutex
}

var consumerSequence atomic.Uint64

func RabbitDial() (*Client, error) {
	rabbiturl := os.Getenv(envvar.RabbitMQURL)
	if strings.TrimSpace(rabbiturl) == "" {
		return nil, errors.New("rabbitmq URL is required")
	}

	conn, err := amqp.Dial(rabbiturl)
	if err != nil {
		return nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil || c.conn.IsClosed() {
		return nil
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("close RabbitMQ connection: %w", err)
	}
	return nil
}

func (c *Client) DeclareQueue(name string, messageTTL time.Duration) (*MessageBroker, error) {
	if c == nil || c.conn == nil || c.conn.IsClosed() {
		return nil, errors.New("RabbitMQ connection unavailable")
	}
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("queue name is required")
	}
	if messageTTL < 0 {
		return nil, errors.New("message TTL cannot be negative")
	}
	if messageTTL > time.Duration(1<<31-1)*time.Millisecond {
		return nil, errors.New("message TTL exceeds RabbitMQ's supported limit")
	}

	channel, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open RabbitMQ channel: %w", err)
	}

	args := amqp.Table{
		amqp.QueueTypeArg: amqp.QueueTypeQuorum,
	}
	if messageTTL > 0 {
		ttlMilliseconds := (messageTTL + time.Millisecond - 1) / time.Millisecond
		args["x-message-ttl"] = int32(ttlMilliseconds)
	}

	if _, err := channel.QueueDeclare(
		name,  //queuename
		true,  //durable
		false, //autoDelete
		false, //exclusive
		false,
		args); err != nil {
		_ = channel.Close()
		return nil, fmt.Errorf("declare RabbitMQ queue %q: %w", name, err)
	}

	return &MessageBroker{queue: name, channel: channel}, nil
}

func (b *MessageBroker) Close() error {
	if b == nil {
		return nil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.closeLocked()
}

func (b *MessageBroker) closeLocked() error {
	if b.channel == nil || b.channel.IsClosed() {
		return nil
	}
	if err := b.channel.Close(); err != nil {
		return fmt.Errorf("close RabbitMQ channel: %w", err)
	}
	return nil
}

func (b *MessageBroker) Publish(ctx context.Context, payload any) error {
	if b == nil {
		return errors.New("RabbitMQ channel unavailable")
	}
	if ctx == nil {
		return errors.New("context is required")
	}
	if payload == nil {
		return errors.New("message payload is required")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	if b.channel == nil || b.channel.IsClosed() {
		return errors.New("RabbitMQ channel unavailable")
	}

	if err := b.channel.PublishWithContext(
		ctx,
		"",
		b.queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("publish message to queue %q: %w", b.queue, err)
	}

	return nil
}

type Handler func(context.Context, []byte) error

func (b *MessageBroker) Consume(ctx context.Context, handler Handler) (returnErr error) {
	if b == nil {
		return errors.New("RabbitMQ channel unavailable")
	}
	if ctx == nil {
		return errors.New("context is required")
	}
	if handler == nil {
		return errors.New("message handler is required")
	}
	b.mu.Lock()
	if b.channel == nil || b.channel.IsClosed() {
		b.mu.Unlock()
		return errors.New("RabbitMQ channel unavailable")
	}
	if err := b.channel.Qos(1, 0, false); err != nil {
		b.mu.Unlock()
		return fmt.Errorf("set RabbitMQ consumer prefetch: %w", err)
	}

	consumerTag := fmt.Sprintf("rabbitmq-%d", consumerSequence.Add(1))
	deliveries, err := b.channel.Consume(
		b.queue,
		consumerTag,
		false,
		false,
		false,
		false,
		nil)
	b.mu.Unlock()
	if err != nil {
		return fmt.Errorf("start consuming queue %q: %w", b.queue, err)
	}
	defer func() {
		b.mu.Lock()
		cancelErr := b.channel.Cancel(consumerTag, false)
		closeErr := b.closeLocked()
		b.mu.Unlock()
		if cancelErr != nil && returnErr == nil {
			returnErr = fmt.Errorf("cancel consumer for queue %q: %w", b.queue, cancelErr)
		}
		if closeErr != nil && returnErr == nil {
			returnErr = closeErr
		}
	}()

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case delivery, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("RabbitMQ consumer for queue %q closed", b.queue)
			}

			if err := handler(ctx, delivery.Body); err != nil {
				b.mu.Lock()
				if nackErr := delivery.Nack(false, true); nackErr != nil {
					b.mu.Unlock()
					return fmt.Errorf("handle message: %v; requeue message: %w", err, nackErr)
				}
				b.mu.Unlock()
				return fmt.Errorf("handle message from queue %q; message requeued: %w", b.queue, err)
			}
			b.mu.Lock()
			if err := delivery.Ack(false); err != nil {
				b.mu.Unlock()
				return fmt.Errorf("acknowledge message from queue %q: %w", b.queue, err)
			}
			b.mu.Unlock()
		}
	}
}

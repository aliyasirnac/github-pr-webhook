package pubsub

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/pkg/config"
	"github.com/aliyasirnac/github-pr-webhook-bot/pkg/pubsubinterface"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"sync"
	"time"
)

type RabbitMQ struct {
	conn     *amqp091.Connection
	channels []*amqp091.Channel
	mu       sync.Mutex
}

func NewRabbitMQ(config *config.PubSubConfig) *RabbitMQ {
	conn, err := amqp091.Dial(config.ConnectionUrl)
	if err != nil {
		zap.L().Fatal("Failed to connect to RabbitMQ", zap.Error(err))
		panic(err)
	}
	return &RabbitMQ{
		conn:     conn,
		channels: []*amqp091.Channel{},
	}
}

func (r *RabbitMQ) Publish(exchange, event string, data []byte) error {
	channel, err := r.conn.Channel()
	if err != nil {
		zap.L().Error("Failed to open a channel", zap.Error(err))
		return err
	}
	defer channel.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		zap.L().Error("Failed to declare exchange", zap.Error(err))
		return err
	}

	return channel.PublishWithContext(ctx, exchange, event, false, false, amqp091.Publishing{
		ContentType:  "text/plain",
		Body:         data,
		DeliveryMode: amqp091.Persistent,
	})
}

func (r *RabbitMQ) Subscribe(exchange, queueName, pattern string, handler func(pubsubinterface.Message)) error {
	channel, err := r.conn.Channel()
	if err != nil {
		zap.L().Error("Failed to open a channel", zap.Error(err))
		return err
	}

	r.mu.Lock()
	r.channels = append(r.channels, channel)
	r.mu.Unlock()

	if err := channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		zap.L().Error("Failed to declare exchange", zap.Error(err))
		return err
	}

	q, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		zap.L().Error("Failed to declare queue", zap.Error(err))
		return err
	}

	if err := channel.QueueBind(q.Name, pattern, exchange, false, nil); err != nil {
		zap.L().Error("Failed to bind queue", zap.Error(err))
		return err
	}

	msgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		zap.L().Error("Failed to register a consumer", zap.Error(err))
		return err
	}

	go func() {
		for msg := range msgs {
			handler(pubsubinterface.Message{
				Topic: exchange,
				Data:  msg.Body,
				Event: pattern,
			})
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, ch := range r.channels {
		_ = ch.Close()
	}
	return r.conn.Close()
}

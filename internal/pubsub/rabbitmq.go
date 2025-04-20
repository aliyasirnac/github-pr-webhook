package pubsub

import (
	"context"
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

func NewRabbitMQ() *RabbitMQ {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		zap.L().Fatal("Failed to connect to RabbitMQ", zap.Error(err))
		panic(err)
	}
	return &RabbitMQ{
		conn:     conn,
		channels: []*amqp091.Channel{},
	}
}

func (r *RabbitMQ) Publish(topic string, data []byte) error {
	channel, err := r.conn.Channel()
	if err != nil {
		zap.L().Error("Failed to open a channel", zap.Error(err))
		return err
	}
	defer channel.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	q, err := channel.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		zap.L().Error("Failed to declare a queue", zap.Error(err))
		return err
	}

	err = channel.Qos(1, 0, false)
	if err != nil {
		zap.L().Error("Failed to set QoS", zap.Error(err))
		return err
	}

	return channel.PublishWithContext(ctx, "", q.Name, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
}

func (r *RabbitMQ) Subscribe(topic string, handler func(pubsubinterface.Message)) error {
	channel, err := r.conn.Channel()
	if err != nil {
		zap.L().Error("Failed to open a channel", zap.Error(err))
		return err
	}

	r.mu.Lock()
	r.channels = append(r.channels, channel)
	r.mu.Unlock()

	q, err := channel.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		zap.L().Error("Failed to declare a queue", zap.Error(err))
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
				Topic: topic,
				Data:  msg.Body,
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

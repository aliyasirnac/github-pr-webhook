package pubsub

import (
	"github.com/aliyasirnac/github-pr-webhook-bot/pkg/pubsubinterface"
	"github.com/nats-io/nats.go"
)

type NatsPubSub struct {
	conn *nats.Conn
}

func NewNatsPubSub(url string) (*NatsPubSub, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsPubSub{conn: conn}, nil
}

func (n *NatsPubSub) Publish(topic, event string, data []byte) error {
	subject := topic + "." + event
	return n.conn.Publish(subject, data)
}

func (n *NatsPubSub) Subscribe(_, queueName, pattern string, h func(pubsubinterface.Message)) error {
	_, err := n.conn.QueueSubscribe(pattern, queueName, func(msg *nats.Msg) {
		h(pubsubinterface.Message{
			Topic: msg.Subject,
			Data:  msg.Data,
			Event: pattern,
		})
	})
	return err
}

func (n *NatsPubSub) Close() error {
	n.conn.Close()
	return nil
}

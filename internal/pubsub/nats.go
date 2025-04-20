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

func (n *NatsPubSub) Publish(topic string, data []byte) error {
	return n.conn.Publish(topic, data)
}

func (n *NatsPubSub) Subscribe(topic string, h func(pubsubinterface.Message)) error {
	_, err := n.conn.Subscribe(topic, func(msg *nats.Msg) {
		h(pubsubinterface.Message{
			Topic: msg.Subject,
			Data:  msg.Data,
		})
	})
	return err
}

func (n *NatsPubSub) Close() error {
	n.conn.Close()
	return nil
}

package pubsubinterface

type Message struct {
	Topic string
	Data  []byte
}

type PubSub interface {
	Publish(topic string, data []byte) error
	Subscribe(topic string, handler func(Message)) error
	Close() error
}

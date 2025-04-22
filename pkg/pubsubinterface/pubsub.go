package pubsubinterface

type Message struct {
	Topic string
	Event string
	Data  []byte
}

type PubSub interface {
	Publish(topic, event string, data []byte) error
	Subscribe(topic, queueName, pattern string, handler func(Message)) error
	Close() error
}

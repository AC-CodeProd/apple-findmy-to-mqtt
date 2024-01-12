package interfaces

type IMQTTClient interface {
	Connect() error
	Disconnect()
	Publish(topic string, payload []byte) error
	Subscribe(topic string, handler MessageHandler) error
}

type MessageHandler func(topic string, payload []byte)

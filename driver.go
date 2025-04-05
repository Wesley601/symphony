package symphony

type QueueDriver interface {
	Subscribe(topic string, handler func(msg []byte)) error
	Publish(topic string, msg []byte) error
}

package drivers

import (
	"github.com/nats-io/nats.go"
)

type NatsDriver struct {
	conn    *nats.Conn
	groupId string
}

func NewNatsDriver(url, groupId string) (*NatsDriver, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsDriver{conn, groupId}, nil
}

func (d *NatsDriver) Subscribe(topic string, handler func(msg []byte)) error {
	_, err := d.conn.QueueSubscribe(topic, d.groupId, func(msg *nats.Msg) {
		handler(msg.Data)
	})

	return err
}

func (d *NatsDriver) Publish(topic string, msg []byte) error {
	return d.conn.Publish(topic, msg)
}

func (d *NatsDriver) Close() {
	d.conn.Drain()
}

package nats

import (
	"github.com/nats-io/nats.go"
)

type Nats struct {
	conn    *nats.Conn
	groupId string
}

func New(url, groupId string) (*Nats, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Nats{conn, groupId}, nil
}

func (d *Nats) Subscribe(topic string, handler func(msg []byte)) error {
	_, err := d.conn.QueueSubscribe(topic, d.groupId, func(msg *nats.Msg) {
		handler(msg.Data)
	})

	return err
}

func (d *Nats) Publish(topic string, msg []byte) error {
	return d.conn.Publish(topic, msg)
}

func (d *Nats) Close() {
	d.conn.Drain()
}

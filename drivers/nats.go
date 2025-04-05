package drivers

import "github.com/nats-io/nats.go"

type NatsDriver struct {
	conn *nats.Conn
}

func NewNatsDriver(url string) (*NatsDriver, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsDriver{conn}, nil
}

func (d *NatsDriver) Subscribe(topic string, handler func(msg []byte)) error {
	_, err := d.conn.Subscribe(topic, func(msg *nats.Msg) {
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

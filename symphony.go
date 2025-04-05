package symphony

import (
	"context"
	"log/slog"

	"github.com/wesley601/symphony/slogutils"
)

type EventHandler interface {
	Handle(context.Context, []byte) ([]byte, error)
}

type Symphony struct {
	queueDriver QueueDriver
	queues      []string
	handlers    []EventHandler
}

func New(queueDriver QueueDriver) *Symphony {
	return &Symphony{
		queueDriver: queueDriver,
		queues:      []string{},
		handlers:    []EventHandler{},
	}
}

func (es *Symphony) After(queue string, handler EventHandler) *Symphony {
	es.queues = append(es.queues, queue)
	es.handlers = append(es.handlers, handler)
	return es
}

func (es *Symphony) Play(ctx context.Context) error {
	for i, queue := range es.queues {
		handler := es.handlers[i]
		slog.Info("Starting Symphony on", slog.String("queue", queue), slogutils.InstanceName(handler))
		es.queueDriver.Subscribe(queue, func(data []byte) {
			data, err := handler.Handle(ctx, data)
			if err != nil {
				slog.Error("Failed to handle event",
					slogutils.Error(err),
					slog.String("queue", queue),
					slogutils.InstanceName(handler),
				)
				return
			}

			if i < len(es.queues)-1 {
				if err := es.queueDriver.Publish(es.queues[i+1], data); err != nil {
					slog.Error("Failed to publish event", slogutils.Error(err))
				}
			}
		})

	}
	return nil
}

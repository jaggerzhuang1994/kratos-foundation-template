package kafka_consumer

import "context"

type Consumer struct {
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Run(_ context.Context) error {
	// todo loop
	return nil
}

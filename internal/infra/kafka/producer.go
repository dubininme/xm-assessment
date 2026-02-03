package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dubininme/xm-assessment/internal/domain/events"
	"github.com/segmentio/kafka-go"
	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
	topic  string
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{}, // Hash by key (companyID) for ordering
	}

	return &Producer{
		writer: writer,
		topic:  topic,
	}
}

func (p *Producer) Publish(ctx context.Context, event events.Event) error {
	payload, err := json.Marshal(event.Payload())
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(event.AggregateID()),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "event_name", Value: []byte(event.EventName())},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (p *Producer) PublishRaw(ctx context.Context, key string, eventType string, payload json.RawMessage) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(payload),
		Headers: []kafka.Header{
			{Key: "event_name", Value: []byte(eventType)},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (p *Producer) PublishBatch(ctx context.Context, messages []kafka.Message) error {
	if err := p.writer.WriteMessages(ctx, messages...); err != nil {
		return fmt.Errorf("failed to write messages to kafka: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

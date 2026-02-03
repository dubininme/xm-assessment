package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/dubininme/xm-assessment/internal/infra/kafka"
	"github.com/dubininme/xm-assessment/internal/infra/postgres"
	"github.com/dubininme/xm-assessment/pkg/logger"
	kafkago "github.com/segmentio/kafka-go"
)

type Processor struct {
	outboxRepo *postgres.OutboxRepo
	producer   *kafka.Producer
	txManager  *postgres.TxManager
	batchSize  int
	interval   time.Duration
}

func NewProcessor(
	outboxRepo *postgres.OutboxRepo,
	producer *kafka.Producer,
	txManager *postgres.TxManager,
	batchSize int,
	interval time.Duration,
) *Processor {
	return &Processor{
		outboxRepo: outboxRepo,
		producer:   producer,
		txManager:  txManager,
		batchSize:  batchSize,
		interval:   interval,
	}
}

func (p *Processor) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	log.Info("outbox processor started",
		"batch_size", p.batchSize,
		"interval", p.interval)

	for {
		select {
		case <-ctx.Done():
			log.Info("outbox processor stopping")
			return p.producer.Close()
		case <-ticker.C:
			if err := p.processBatch(ctx); err != nil {
				log.Error("error processing outbox batch", "error", err)
			}
		}
	}
}

func (p *Processor) processBatch(ctx context.Context) error {
	log := logger.FromContext(ctx)
	var events []postgres.OutboxEvent
	var messages []kafkago.Message

	err := p.txManager.Do(ctx, func(txCtx context.Context) error {
		var err error
		events, err = p.outboxRepo.GetUnprocessed(txCtx, p.batchSize)
		if err != nil {
			return fmt.Errorf("failed to get unprocessed events: %w", err)
		}

		if len(events) == 0 {
			return nil
		}

		ids := make([]int64, len(events))
		for i, e := range events {
			ids[i] = e.ID
		}

		if err := p.outboxRepo.MarkProcessed(txCtx, ids); err != nil {
			return fmt.Errorf("failed to mark events as processed: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}

	messages = make([]kafkago.Message, 0, len(events))
	for _, e := range events {
		messages = append(messages, kafkago.Message{
			Key:   []byte(e.AggregateID),
			Value: []byte(e.Payload),
			Headers: []kafkago.Header{
				{Key: "event_name", Value: []byte(e.EventType)},
			},
		})
	}

	if err := p.producer.PublishBatch(ctx, messages); err != nil {
		log.Warn("failed to send events to Kafka (events marked as processed)",
			"count", len(events),
			"error", err)
		return fmt.Errorf("failed to write messages to kafka: %w", err)
	}

	log.Info("processed events from outbox", "count", len(events))
	return nil
}

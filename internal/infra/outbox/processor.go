package outbox

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dubininme/xm-assessment/internal/infra/kafka"
	"github.com/dubininme/xm-assessment/internal/infra/postgres"
	"github.com/dubininme/xm-assessment/pkg/logger"
	kafkago "github.com/segmentio/kafka-go"
)

type Processor struct {
	outboxRepo     *postgres.OutboxRepo
	producer       *kafka.Producer
	txManager      *postgres.TxManager
	batchSize      int
	interval       time.Duration
	publishTimeout time.Duration
}

func NewProcessor(
	outboxRepo *postgres.OutboxRepo,
	producer *kafka.Producer,
	txManager *postgres.TxManager,
	batchSize int,
	interval time.Duration,
	publishTimeout time.Duration,
) *Processor {
	return &Processor{
		outboxRepo:     outboxRepo,
		producer:       producer,
		txManager:      txManager,
		batchSize:      batchSize,
		interval:       interval,
		publishTimeout: publishTimeout,
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

	return p.txManager.Do(ctx, func(txCtx context.Context) error {
		events, err := p.outboxRepo.GetUnprocessed(txCtx, p.batchSize)
		if err != nil {
			return fmt.Errorf("failed to get unprocessed events: %w", err)
		}

		if len(events) == 0 {
			return nil
		}

		messages := make([]kafkago.Message, 0, len(events))
		ids := make([]int64, 0, len(events))
		for _, e := range events {
			messages = append(messages, kafkago.Message{
				Key:   []byte(e.AggregateID),
				Value: []byte(e.Payload),
				Headers: []kafkago.Header{
					{Key: "event_name", Value: []byte(e.EventType)},
					{Key: "outbox_id", Value: []byte(strconv.FormatInt(e.ID, 10))},
				},
			})
			ids = append(ids, e.ID)
		}

		publishCtx, cancel := context.WithTimeout(txCtx, p.publishTimeout)
		defer cancel()

		publishErr := p.producer.PublishBatch(publishCtx, messages)
		if publishErr != nil {
			return fmt.Errorf("failed to write messages to kafka: %w", publishErr)
		}

		markErr := p.outboxRepo.MarkProcessed(txCtx, ids)
		if markErr != nil {
			return fmt.Errorf("failed to mark events as processed: %w", markErr)
		}

		log.Info("processed events from outbox", "count", len(events))
		return nil
	})
}

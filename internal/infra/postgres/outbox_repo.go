package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dubininme/xm-assessment/internal/domain/events"
	"github.com/lib/pq"
)

var _ events.EventsPublisher = (*OutboxRepo)(nil)

type OutboxRepo struct {
	db *Db
}

func NewOutboxRepo(db *Db) *OutboxRepo {
	return &OutboxRepo{db: db}
}

func (r *OutboxRepo) Publish(ctx context.Context, event events.Event) error {
	exec := ExtractExecutor(ctx, r.db)
	query := `INSERT INTO outbox (event_type, aggregate_id, payload, created_at)
	          VALUES ($1, $2, $3, $4)`

	payload, err := json.Marshal(event.Payload())
	if err != nil {
		return err
	}

	_, err = exec.ExecContext(ctx, query,
		event.EventName(),
		event.AggregateID(),
		payload,
		event.CreatedAt(),
	)
	return err
}

func (r *OutboxRepo) MarkProcessed(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	exec := ExtractExecutor(ctx, r.db)
	query := `UPDATE outbox SET is_processed = true, processed_at = $1 WHERE id = ANY($2)`
	_, err := exec.ExecContext(ctx, query, time.Now().Unix(), pq.Array(ids))
	return err
}

func (r *OutboxRepo) GetUnprocessed(ctx context.Context, limit int) ([]OutboxEvent, error) {
	exec := ExtractExecutor(ctx, r.db)
	query := `SELECT id, event_type, aggregate_id, payload, created_at
	          FROM outbox
	          WHERE is_processed = false
	          ORDER BY id ASC
	          LIMIT $1
	          FOR UPDATE SKIP LOCKED`

	rows, err := exec.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var events []OutboxEvent
	for rows.Next() {
		var e OutboxEvent
		if err := rows.Scan(&e.ID, &e.EventType, &e.AggregateID, &e.Payload, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, rows.Err()
}

type OutboxEvent struct {
	ID          int64
	EventType   string
	AggregateID string
	Payload     json.RawMessage
	CreatedAt   int64
}

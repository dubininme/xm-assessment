CREATE TABLE outbox (
    id BIGSERIAL PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    aggregate_id UUID NOT NULL,
    payload JSONB NOT NULL,
    created_at BIGINT NOT NULL,
    processed_at BIGINT,
    is_processed BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_outbox_unprocessed ON outbox(created_at) WHERE is_processed = false;
CREATE INDEX idx_outbox_aggregate_id ON outbox(aggregate_id);
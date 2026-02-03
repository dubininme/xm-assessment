package handler

import "context"

type HealthChecker interface {
	Check(ctx context.Context) error
	Name() string
}

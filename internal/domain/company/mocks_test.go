package company

import (
	"context"

	"github.com/dubininme/xm-assessment/internal/domain/events"
	"github.com/stretchr/testify/mock"
)

// MockCompanyRepository is a mock implementation of CompanyRepository interface
type MockCompanyRepository struct {
	mock.Mock
}

func (m *MockCompanyRepository) Create(ctx context.Context, c Company) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCompanyRepository) GetByID(ctx context.Context, id string) (*Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Company), args.Error(1)
}

func (m *MockCompanyRepository) Update(ctx context.Context, c Company) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockCompanyRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEventsPublisher is a mock implementation of EventsPublisher interface
type MockEventsPublisher struct {
	mock.Mock
}

func (m *MockEventsPublisher) Publish(ctx context.Context, event events.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// MockTxManager is a mock implementation of TxManager interface
type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)

	// Always execute the function to simulate real transaction behavior
	fnErr := fn(ctx)

	// If function returned error, return it (simulates rollback)
	if fnErr != nil {
		return fnErr
	}

	// If mock specifies commit error, return it
	return args.Error(0)
}

func isCompanyCreatedEvent() interface{} {
	return mock.MatchedBy(func(e events.Event) bool {
		_, ok := e.(CompanyCreatedEvent)
		if !ok {
			// Try pointer type
			_, ok = e.(*CompanyCreatedEvent)
		}
		return ok
	})
}

func isCompanyUpdatedEvent() interface{} {
	return mock.MatchedBy(func(e events.Event) bool {
		_, ok := e.(CompanyUpdatedEvent)
		if !ok {
			_, ok = e.(*CompanyUpdatedEvent)
		}
		return ok
	})
}

func isCompanyDeletedEvent() interface{} {
	return mock.MatchedBy(func(e events.Event) bool {
		_, ok := e.(CompanyDeletedEvent)
		if !ok {
			_, ok = e.(*CompanyDeletedEvent)
		}
		return ok
	})
}

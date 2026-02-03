package company

import (
	"context"
	"fmt"

	"github.com/dubininme/xm-assessment/internal/domain/events"
	"github.com/google/uuid"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type CompanyService struct {
	repo      CompanyRepository
	publisher events.EventsPublisher
	txManager TxManager
}

func NewCompanyService(repo CompanyRepository, publisher events.EventsPublisher, txManager TxManager) *CompanyService {
	return &CompanyService{repo: repo, publisher: publisher, txManager: txManager}
}

func (s *CompanyService) CreateCompany(ctx context.Context, params CreateParams) (*Company, error) {
	c, err := NewCompany(uuid.New(), params.Name, params.Description, params.EmployeesCount, params.Type)
	if err != nil {
		return nil, err
	}

	// Set registered status from params
	if params.Registered {
		c.Register()
	}

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		err := s.repo.Create(ctx, *c)
		if err != nil {
			return fmt.Errorf("failed to create company: %w", err)
		}

		err = s.publisher.Publish(ctx, NewCompanyCreatedEvent(
			c.ID().String(),
			params,
		))

		if err != nil {
			return fmt.Errorf("failed to publish company created event: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, companyID string, params UpdateParams) (*Company, error) {
	c, err := s.repo.GetByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company by ID: %w", err)
	}

	if params.Name != nil {
		err := c.SetName(*params.Name)
		if err != nil {
			// Return validation errors directly without wrapping
			return nil, err
		}
	}

	if params.Description != nil {
		err := c.SetDescription(*params.Description)
		if err != nil {
			// Return validation errors directly without wrapping
			return nil, err
		}
	}

	if params.EmployeesCount != nil {
		err := c.SetEmployeesCount(*params.EmployeesCount)
		if err != nil {
			// Return validation errors directly without wrapping
			return nil, err
		}
	}

	if params.Type != nil {
		err := c.SetType(*params.Type)
		if err != nil {
			// Return validation errors directly without wrapping
			return nil, err
		}
	}

	if params.Registered != nil {
		c.SetRegistered(*params.Registered)
	}

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		err := s.repo.Update(ctx, *c)
		if err != nil {
			return fmt.Errorf("failed to update company: %w", err)
		}

		err = s.publisher.Publish(ctx, NewCompanyUpdatedEvent(
			c.ID().String(),
			params,
		))

		if err != nil {
			return fmt.Errorf("failed to publish company updated event: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *CompanyService) GetByID(ctx context.Context, companyID string) (*Company, error) {
	company, err := s.repo.GetByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company by ID: %w", err)
	}

	return company, nil
}

func (s *CompanyService) DeleteCompany(ctx context.Context, companyID string) error {
	return s.txManager.Do(ctx, func(ctx context.Context) error {
		err := s.repo.Delete(ctx, companyID)
		if err != nil {
			return fmt.Errorf("failed to delete company: %w", err)
		}

		err = s.publisher.Publish(ctx, NewCompanyDeletedEvent(companyID))
		if err != nil {
			return fmt.Errorf("failed to publish company deleted event: %w", err)
		}

		return nil
	})
}

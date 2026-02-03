package company

import (
	"context"
)

type CompanyRepository interface {
	Create(ctx context.Context, company Company) error
	Update(ctx context.Context, company Company) error
	Delete(ctx context.Context, companyID string) error
	GetByID(ctx context.Context, companyID string) (*Company, error)
}

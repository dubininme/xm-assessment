package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dubininme/xm-assessment/internal/domain/company"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const ErrUniqueViolationCode = "23505"

type CompanyRepo struct {
	db *Db
}

var _ company.CompanyRepository = (*CompanyRepo)(nil)

func NewCompanyRepo(db *Db) *CompanyRepo {
	return &CompanyRepo{db: db}
}

func (r *CompanyRepo) Create(ctx context.Context, c company.Company) error {
	exec := ExtractExecutor(ctx, r.db)
	_, err := exec.ExecContext(ctx, `
		INSERT INTO companies (id, name, description, employees_count, registered, type)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		c.ID().String(), c.Name().String(), c.Description().String(), c.EmployeesCount().Int(), c.IsRegistered(), c.CompanyType().Int())

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == ErrUniqueViolationCode {
				return company.ErrCompanyNameAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (r *CompanyRepo) Update(ctx context.Context, c company.Company) error {
	exec := ExtractExecutor(ctx, r.db)
	res, err := exec.ExecContext(ctx, `
		UPDATE companies SET name = $2, description = $3, employees_count = $4, registered = $5, type = $6
		WHERE id = $1`,
		c.ID().String(), c.Name().String(), c.Description().String(), c.EmployeesCount().Int(), c.IsRegistered(), c.CompanyType().Int(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == ErrUniqueViolationCode {
				return company.ErrCompanyNameAlreadyExists
			}
		}

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return company.ErrCompanyNotFound
	}

	return nil
}

func (r *CompanyRepo) GetByID(ctx context.Context, companyID string) (*company.Company, error) {
	exec := ExtractExecutor(ctx, r.db)
	row := exec.QueryRowContext(ctx, `
		SELECT id, name, description, employees_count, registered, type
		FROM companies WHERE id = $1`, companyID)

	var queryResult CompanyRowDto
	err := row.Scan(&queryResult.ID, &queryResult.Name, &queryResult.Description, &queryResult.EmployeesCount, &queryResult.Registered, &queryResult.Type)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, company.ErrCompanyNotFound
		}

		return nil, err
	}

	return queryResult.ToEntity()
}

func (r *CompanyRepo) Delete(ctx context.Context, companyID string) error {
	exec := ExtractExecutor(ctx, r.db)
	res, err := exec.ExecContext(ctx, `DELETE FROM companies WHERE id = $1`, companyID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return company.ErrCompanyNotFound
	}

	return nil
}

type CompanyRowDto struct {
	ID             string
	Name           string
	Description    string
	EmployeesCount int
	Registered     bool
	Type           int16
}

func (r *CompanyRowDto) ToEntity() (*company.Company, error) {
	cType, err := company.CompanyTypeFromInt(r.Type)
	if err != nil {
		return nil, err
	}

	c, err := company.NewCompany(
		uuid.MustParse(r.ID),
		r.Name,
		r.Description,
		r.EmployeesCount,
		cType.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid data from database: %w", err)
	}

	if r.Registered {
		c.Register()
	}

	return c, nil
}

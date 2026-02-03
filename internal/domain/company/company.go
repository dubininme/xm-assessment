package company

import (
	"fmt"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	CorporationsType       CompanyType = "Corporations"
	NonProfitType          CompanyType = "NonProfit"
	CooperativeType        CompanyType = "Cooperative"
	SoleProprietorshipType CompanyType = "Sole Proprietorship"

	CorporationsIntType       = 1
	NonProfitIntType          = 2
	CooperativeIntType        = 3
	SoleProprietorshipIntType = 4
)

type CompanyType string

func (t CompanyType) String() string {
	return string(t)
}

func (t CompanyType) Int() int16 {
	switch t {
	case CorporationsType:
		return CorporationsIntType
	case NonProfitType:
		return NonProfitIntType
	case CooperativeType:
		return CooperativeIntType
	case SoleProprietorshipType:
		return SoleProprietorshipIntType
	default:
		return 0
	}
}

func CompanyTypeFromInt(i int16) (CompanyType, error) {
	switch i {
	case CorporationsIntType:
		return CorporationsType, nil
	case NonProfitIntType:
		return NonProfitType, nil
	case CooperativeIntType:
		return CooperativeType, nil
	case SoleProprietorshipIntType:
		return SoleProprietorshipType, nil
	default:
		return "", ErrInvalidCompanyType
	}
}

func NewCompanyType(s string) (*CompanyType, error) {
	cType := CompanyType(s)
	switch cType {
	case CorporationsType, NonProfitType, CooperativeType, SoleProprietorshipType:
		return &cType, nil
	default:
		return nil, ErrInvalidCompanyType
	}
}

type CompanyName string

func (n CompanyName) String() string {
	return string(n)
}

func NewCompanyName(name string) (*CompanyName, error) {
	len := utf8.RuneCountInString(name)
	if len == 0 || len > 15 {
		return nil, ErrInvalidCompanyNameLength
	}

	n := CompanyName(name)
	return &n, nil
}

type CompanyDescription string

func (d CompanyDescription) String() string {
	return string(d)
}

func NewCompanyDescription(description string) (*CompanyDescription, error) {
	len := utf8.RuneCountInString(description)

	if len > 3000 {
		return nil, ErrInvalidCompanyDescriptionLength
	}

	d := CompanyDescription(description)
	return &d, nil
}

type EmployeesCount int

func (c EmployeesCount) Int() int {
	return int(c)
}

func NewEmployeesCount(count int) (*EmployeesCount, error) {
	if count < 1 {
		return nil, ErrInvalidEmployeesCount
	}

	c := EmployeesCount(count)
	return &c, nil
}

type Company struct {
	id             uuid.UUID
	name           CompanyName
	description    CompanyDescription
	employeesCount EmployeesCount
	registered     bool
	cType          CompanyType
}

func NewCompany(id uuid.UUID, name string, description string, employeesCount int, companyType string) (*Company, error) {
	cName, err := NewCompanyName(name)
	if err != nil {
		return nil, fmt.Errorf("error creating company name: %w", err)
	}

	cDescription, err := NewCompanyDescription(description)
	if err != nil {
		return nil, fmt.Errorf("error creating company description: %w", err)
	}

	eCount, err := NewEmployeesCount(employeesCount)
	if err != nil {
		return nil, fmt.Errorf("error creating employees count: %w", err)
	}

	cType, err := NewCompanyType(companyType)
	if err != nil {
		return nil, fmt.Errorf("error creating company type: %w", err)
	}

	return &Company{
		id:             id,
		name:           *cName,
		description:    *cDescription,
		employeesCount: *eCount,
		registered:     false,
		cType:          *cType,
	}, nil

}

func (c *Company) ID() uuid.UUID {
	return c.id
}

func (c *Company) Name() CompanyName {
	return c.name
}

func (c *Company) Description() CompanyDescription {
	return c.description
}

func (c *Company) EmployeesCount() EmployeesCount {
	return c.employeesCount
}

func (c *Company) CompanyType() CompanyType {
	return c.cType
}

func (c *Company) IsRegistered() bool {
	return c.registered
}

func (c *Company) Register() {
	c.registered = true
}

func (c *Company) SetName(name string) error {
	n, err := NewCompanyName(name)
	if err != nil {
		return err
	}
	c.name = *n
	return nil
}

func (c *Company) SetDescription(desc string) error {
	d, err := NewCompanyDescription(desc)
	if err != nil {
		return err
	}
	c.description = *d
	return nil
}

func (c *Company) SetEmployeesCount(count int) error {
	e, err := NewEmployeesCount(count)
	if err != nil {
		return err
	}
	c.employeesCount = *e
	return nil
}

func (c *Company) SetType(t string) error {
	ct, err := NewCompanyType(t)
	if err != nil {
		return err
	}
	c.cType = *ct
	return nil
}

func (c *Company) SetRegistered(r bool) {
	c.registered = r
}

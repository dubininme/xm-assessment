package handler

import (
	"testing"

	"github.com/dubininme/xm-assessment/internal/domain/company"
	"github.com/dubininme/xm-assessment/pkg/gen/oapi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompanyToResponse(t *testing.T) {
	id := uuid.New()
	c, err := company.NewCompany(id, "TestCo", "Test Description", 100, "Corporations")
	require.NoError(t, err)

	response := CompanyToResponse(c)

	assert.Equal(t, id, uuid.UUID(response.Id))
	assert.Equal(t, "TestCo", response.Name)
	assert.Equal(t, 100, response.EmployeesCount)
	assert.False(t, response.Registered)
	assert.Equal(t, oapi.CompanyType("Corporations"), response.Type)
}

func TestCreateRequestToParams(t *testing.T) {
	desc := "Test Description"
	req := oapi.CreateCompanyRequest{
		Name:           "NewCo",
		Description:    &desc,
		EmployeesCount: 50,
		Registered:     true,
		Type:           oapi.CompanyType("NonProfit"),
	}

	params := CreateRequestToParams(req)

	assert.Equal(t, "NewCo", params.Name)
	assert.Equal(t, "Test Description", params.Description)
	assert.Equal(t, 50, params.EmployeesCount)
	assert.True(t, params.Registered)
	assert.Equal(t, "NonProfit", params.Type)
}

func TestUpdateRequestToParams_AllFields(t *testing.T) {
	name := "Updated"
	desc := "Updated Desc"
	count := 200
	registered := true
	companyType := oapi.CompanyType("SoleProprietorship")

	req := oapi.UpdateCompanyRequest{
		Name:           &name,
		Description:    &desc,
		EmployeesCount: &count,
		Registered:     &registered,
		Type:           &companyType,
	}

	params := UpdateRequestToParams(req)

	assert.Equal(t, &name, params.Name)
	assert.Equal(t, &desc, params.Description)
	assert.Equal(t, &count, params.EmployeesCount)
	assert.Equal(t, &registered, params.Registered)
	require.NotNil(t, params.Type)
	assert.Equal(t, "SoleProprietorship", *params.Type)
}

func TestUpdateRequestToParams_PartialFields(t *testing.T) {
	name := "OnlyName"
	req := oapi.UpdateCompanyRequest{
		Name: &name,
	}

	params := UpdateRequestToParams(req)

	assert.Equal(t, &name, params.Name)
	assert.Nil(t, params.Description)
	assert.Nil(t, params.EmployeesCount)
	assert.Nil(t, params.Registered)
	assert.Nil(t, params.Type)
}

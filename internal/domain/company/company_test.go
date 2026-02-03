package company

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCompany_Success(t *testing.T) {
	id := uuid.New()
	c, err := NewCompany(id, "ValidName", "Description", 10, "Corporations")

	require.NoError(t, err)
	require.NotNil(t, c)
	assert.Equal(t, id, c.ID())
	assert.Equal(t, "ValidName", c.Name().String())
	assert.Equal(t, "Description", c.Description().String())
	assert.Equal(t, 10, c.EmployeesCount().Int())
	assert.Equal(t, CorporationsType, c.CompanyType())
	assert.False(t, c.IsRegistered()) // default is false
}

func TestNewCompany_InvalidName(t *testing.T) {
	id := uuid.New()
	c, err := NewCompany(id, "", "Description", 10, "Corporations")

	require.Error(t, err)
	require.Nil(t, c)
	assert.Contains(t, err.Error(), "company name")
}

func TestNewCompany_InvalidEmployeesCount(t *testing.T) {
	id := uuid.New()
	c, err := NewCompany(id, "ValidName", "Description", 0, "Corporations")

	require.Error(t, err)
	require.Nil(t, c)
	assert.Contains(t, err.Error(), "employees count")
}

func TestCompany_SetName_Success(t *testing.T) {
	c, _ := NewCompany(uuid.New(), "Original", "Desc", 10, "Corporations")

	err := c.SetName("NewName")
	require.NoError(t, err)
	assert.Equal(t, "NewName", c.Name().String())
}

func TestCompany_SetName_Invalid(t *testing.T) {
	c, _ := NewCompany(uuid.New(), "Original", "Desc", 10, "Corporations")

	err := c.SetName("")
	require.ErrorIs(t, err, ErrInvalidCompanyNameLength)
}

func TestCompany_SetEmployeesCount_Success(t *testing.T) {
	c, _ := NewCompany(uuid.New(), "Name", "Desc", 10, "Corporations")

	err := c.SetEmployeesCount(100)
	require.NoError(t, err)
	assert.Equal(t, 100, c.EmployeesCount().Int())
}

func TestCompany_SetEmployeesCount_Invalid(t *testing.T) {
	c, _ := NewCompany(uuid.New(), "Name", "Desc", 10, "Corporations")

	err := c.SetEmployeesCount(0)
	require.ErrorIs(t, err, ErrInvalidEmployeesCount)
}

func TestCompany_Register(t *testing.T) {
	c, _ := NewCompany(uuid.New(), "Name", "Desc", 10, "Corporations")

	assert.False(t, c.IsRegistered())

	c.Register()
	assert.True(t, c.IsRegistered())
}

func TestCompany_SetRegistered(t *testing.T) {
	c, _ := NewCompany(uuid.New(), "Name", "Desc", 10, "Corporations")

	c.Register()
	assert.True(t, c.IsRegistered())

	c.SetRegistered(false)
	assert.False(t, c.IsRegistered())
}

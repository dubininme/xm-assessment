// +build unit

package company

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupServiceMocks(t *testing.T) (*CompanyService, *MockCompanyRepository, *MockEventsPublisher, *MockTxManager) {
	mockRepo := new(MockCompanyRepository)
	mockPublisher := new(MockEventsPublisher)
	mockTxManager := new(MockTxManager)
	service := NewCompanyService(mockRepo, mockPublisher, mockTxManager)

	return service, mockRepo, mockPublisher, mockTxManager
}

func TestCreateCompany_Success(t *testing.T) {
	service, mockRepo, mockPublisher, mockTxManager := setupServiceMocks(t)

	params := CreateParams{
		Name:           "TechCorp",
		Description:    "A tech company",
		EmployeesCount: 50,
		Registered:     false,
		Type:           "Corporations",
	}

	ctx := context.Background()

	mockTxManager.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("company.Company")).Return(nil)
	mockPublisher.On("Publish", mock.Anything, isCompanyCreatedEvent()).Return(nil)

	company, err := service.CreateCompany(ctx, params)

	require.NoError(t, err)
	require.NotNil(t, company)
	assert.Equal(t, "TechCorp", company.Name().String())
	assert.Equal(t, 50, company.EmployeesCount().Int())
	assert.False(t, company.IsRegistered())

	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
	mockTxManager.AssertExpectations(t)
}

func TestCreateCompany_EmptyName(t *testing.T) {
	service, mockRepo, mockPublisher, _ := setupServiceMocks(t)

	params := CreateParams{
		Name:           "",
		EmployeesCount: 10,
		Type:           "Corporations",
	}

	result, err := service.CreateCompany(context.Background(), params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrInvalidCompanyNameLength)

	mockRepo.AssertNotCalled(t, "Create")
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestCreateCompany_RepoError(t *testing.T) {
	service, mockRepo, _, mockTxManager := setupServiceMocks(t)

	params := CreateParams{
		Name:           "TestCo",
		EmployeesCount: 10,
		Type:           "Corporations",
	}

	repoErr := errors.New("database error")
	mockTxManager.On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(repoErr)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("Company")).Return(repoErr)

	result, err := service.CreateCompany(context.Background(), params)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateCompany_Success(t *testing.T) {
	service, mockRepo, mockPublisher, mockTxManager := setupServiceMocks(t)

	companyID := uuid.New()
	existingCompany, _ := NewCompany(companyID, "OldName", "Old Desc", 5, "Corporations")

	newName := "NewName"
	newDesc := "New Description"
	params := UpdateParams{
		Name:        &newName,
		Description: &newDesc,
	}

	mockRepo.On("GetByID", mock.Anything, companyID.String()).Return(existingCompany, nil)
	mockTxManager.On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("Company")).Return(nil)
	mockPublisher.On("Publish", mock.Anything, isCompanyUpdatedEvent()).Return(nil)

	result, err := service.UpdateCompany(context.Background(), companyID.String(), params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "NewName", result.Name().String())
	assert.Equal(t, "New Description", result.Description().String())
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestUpdateCompany_NotFound(t *testing.T) {
	service, mockRepo, mockPublisher, _ := setupServiceMocks(t)

	companyID := uuid.New().String()
	newName := "NewName"
	params := UpdateParams{Name: &newName}

	mockRepo.On("GetByID", mock.Anything, companyID).Return(nil, ErrCompanyNotFound)

	result, err := service.UpdateCompany(context.Background(), companyID, params)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestDeleteCompany_Success(t *testing.T) {
	service, mockRepo, mockPublisher, mockTxManager := setupServiceMocks(t)

	companyID := uuid.New().String()

	mockTxManager.On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockRepo.On("Delete", mock.Anything, companyID).Return(nil)
	mockPublisher.On("Publish", mock.Anything, isCompanyDeletedEvent()).Return(nil)

	err := service.DeleteCompany(context.Background(), companyID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestDeleteCompany_NotFound(t *testing.T) {
	service, mockRepo, mockPublisher, mockTxManager := setupServiceMocks(t)

	companyID := uuid.New().String()

	mockTxManager.On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(ErrCompanyNotFound)
	mockRepo.On("Delete", mock.Anything, companyID).Return(ErrCompanyNotFound)

	err := service.DeleteCompany(context.Background(), companyID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrCompanyNotFound)
	mockRepo.AssertExpectations(t)
	// Важно: событие НЕ должно быть опубликовано
	mockPublisher.AssertNotCalled(t, "Publish")
}

func TestDeleteCompany_RepoError(t *testing.T) {
	service, mockRepo, _, mockTxManager := setupServiceMocks(t)

	companyID := uuid.New().String()
	repoErr := errors.New("delete failed")

	mockTxManager.On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockRepo.On("Delete", mock.Anything, companyID).Return(repoErr)

	err := service.DeleteCompany(context.Background(), companyID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

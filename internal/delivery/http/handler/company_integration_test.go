//go:build integration

package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dubininme/xm-assessment/internal/config"
	deliveryHttp "github.com/dubininme/xm-assessment/internal/delivery/http"
	"github.com/dubininme/xm-assessment/internal/delivery/http/handler"
	"github.com/dubininme/xm-assessment/internal/delivery/http/middleware"
	"github.com/dubininme/xm-assessment/internal/domain/company"
	"github.com/dubininme/xm-assessment/internal/infra/auth"
	"github.com/dubininme/xm-assessment/internal/infra/postgres"
	"github.com/dubininme/xm-assessment/pkg/gen/oapi"
	"github.com/dubininme/xm-assessment/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration test that requires a running PostgreSQL database
// Run with: docker-compose up -d postgres && go test -v ./internal/delivery/http/handler/...
func TestCompanyLifecycle_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	ctx = logger.WithLogger(ctx, logger.NewTestLogger())

	// Setup database connection (using docker-compose postgres)
	cfg := &config.DbConfig{
		DBHost:            "localhost",
		DBPort:            "5432",
		DBUser:            "xm_user",
		DBPassword:        "xm_password",
		DBName:            "xm_db",
		DBMaxOpenConns:    10,
		DBMaxIdleConns:    5,
		DBConnMaxLifetime: 300,
	}

	db, err := postgres.Connect(ctx, *cfg)
	require.NoError(t, err, "Failed to connect to test database. Make sure docker-compose is running.")
	defer func() { _ = db.Close() }()

	_, _ = db.ExecContext(ctx, "DELETE FROM companies WHERE name LIKE 'IntegrationTest%'")

	companyRepo := postgres.NewCompanyRepo(db)
	outboxRepo := postgres.NewOutboxRepo(db)
	txManager := postgres.NewTxManager(db)
	dbChecker := postgres.NewDBHealthChecker(db)

	companyService := company.NewCompanyService(companyRepo, outboxRepo, txManager)
	companyHandler := handler.NewCompanyHandler(companyService)
	healthHandler := handler.NewHealthHandler(dbChecker)

	jwtService := auth.NewJWTService("test-secret-key-for-integration-tests")
	authHandler := handler.NewAuthHandler(jwtService)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	router := deliveryHttp.NewRouter(companyHandler, healthHandler, authHandler, authMiddleware)

	token := getAuthToken(t, router)

	t.Log("Test 1: Creating company...")
	createReq := oapi.CreateCompanyRequest{
		Name:           "IntegrationTest",
		Description:    ptr("Test company for integration testing"),
		EmployeesCount: 100,
		Registered:     false,
		Type:           oapi.Corporations,
	}

	companyID := createCompany(t, router, token, createReq)
	assert.NotEmpty(t, companyID, "Company ID should not be empty")
	t.Logf("Company created with ID: %s", companyID)

	t.Log("Test 2: Getting company...")
	fetchedCompany := getCompany(t, router, companyID)
	assert.Equal(t, "IntegrationTest", fetchedCompany.Name)
	assert.Equal(t, 100, fetchedCompany.EmployeesCount)
	assert.False(t, fetchedCompany.Registered)
	assert.NotNil(t, fetchedCompany.Description)
	assert.Equal(t, "Test company for integration testing", *fetchedCompany.Description)
	t.Log("Company fetched successfully")

	t.Log("Test 3: Updating company...")
	updateReq := oapi.UpdateCompanyRequest{
		EmployeesCount: ptr(250),
		Registered:     ptr(true),
	}

	updateCompany(t, router, token, companyID, updateReq)

	updatedCompany := getCompany(t, router, companyID)
	assert.Equal(t, 250, updatedCompany.EmployeesCount)
	assert.True(t, updatedCompany.Registered)
	t.Log("Company updated successfully")

	t.Log("Test 4: Deleting company...")
	deleteCompany(t, router, token, companyID)
	t.Log("Company deleted successfully")

	resp := makeRequest(t, router, "GET", fmt.Sprintf("/api/v1/companies/%s", companyID), "", nil)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	t.Log("Verified company is deleted (404)")

	t.Log("Test 5: Testing double delete (should return 404)...")
	resp = makeRequest(t, router, "DELETE", fmt.Sprintf("/api/v1/companies/%s", companyID), token, nil)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	t.Log("Double delete correctly returns 404")
}

func getAuthToken(t *testing.T, router http.Handler) string {
	t.Helper()
	req := map[string]string{"user_id": "test-user"}
	body, _ := json.Marshal(req)
	resp := makeRequest(t, router, "POST", "/api/v1/auth/token", "", body)
	require.Equal(t, http.StatusOK, resp.Code)

	var tokenResp struct {
		Token string `json:"token"`
	}
	err := json.NewDecoder(resp.Body).Decode(&tokenResp)
	require.NoError(t, err)
	return tokenResp.Token
}

func createCompany(t *testing.T, router http.Handler, token string, req oapi.CreateCompanyRequest) string {
	t.Helper()
	body, _ := json.Marshal(req)
	resp := makeRequest(t, router, "POST", "/api/v1/companies", token, body)
	require.Equal(t, http.StatusCreated, resp.Code)

	var company oapi.Company
	err := json.NewDecoder(resp.Body).Decode(&company)
	require.NoError(t, err)
	return company.Id.String()
}

func getCompany(t *testing.T, router http.Handler, id string) oapi.Company {
	t.Helper()
	resp := makeRequest(t, router, "GET", fmt.Sprintf("/api/v1/companies/%s", id), "", nil)
	require.Equal(t, http.StatusOK, resp.Code)

	var company oapi.Company
	err := json.NewDecoder(resp.Body).Decode(&company)
	require.NoError(t, err)
	return company
}

func updateCompany(t *testing.T, router http.Handler, token, id string, req oapi.UpdateCompanyRequest) {
	t.Helper()
	body, _ := json.Marshal(req)
	resp := makeRequest(t, router, "PATCH", fmt.Sprintf("/api/v1/companies/%s", id), token, body)
	require.Equal(t, http.StatusOK, resp.Code)
}

func deleteCompany(t *testing.T, router http.Handler, token, id string) {
	t.Helper()
	resp := makeRequest(t, router, "DELETE", fmt.Sprintf("/api/v1/companies/%s", id), token, nil)
	require.Equal(t, http.StatusNoContent, resp.Code)
}

func makeRequest(t *testing.T, router http.Handler, method, path, token string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func ptr[T any](v T) *T {
	return &v
}

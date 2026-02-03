package http

import (
	"net/http"

	"github.com/dubininme/xm-assessment/internal/delivery/http/handler"
	"github.com/dubininme/xm-assessment/internal/delivery/http/middleware"
	"github.com/gorilla/mux"
)

func NewRouter(
	companyHandler *handler.CompanyHandler,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Health check without prefix (for load balancers, k8s probes, etc.)
	router.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)

	// API v1 routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Public routes
	apiV1.HandleFunc("/companies/{id}", companyHandler.GetCompany).Methods(http.MethodGet)
	apiV1.HandleFunc("/auth/token", authHandler.GenerateToken).Methods(http.MethodPost)

	// Protected routes
	protected := apiV1.NewRoute().Subrouter()
	protected.Use(authMiddleware.Authenticate)

	protected.HandleFunc("/companies", companyHandler.CreateCompany).Methods(http.MethodPost)
	protected.HandleFunc("/companies/{id}", companyHandler.UpdateCompany).Methods(http.MethodPatch)
	protected.HandleFunc("/companies/{id}", companyHandler.DeleteCompany).Methods(http.MethodDelete)

	return router
}

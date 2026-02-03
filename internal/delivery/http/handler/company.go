package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dubininme/xm-assessment/internal/domain/company"
	"github.com/dubininme/xm-assessment/pkg/gen/oapi"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CompanyHandler struct {
	service *company.CompanyService
}

func NewCompanyHandler(service *company.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		writeErr(w, http.StatusBadRequest, oapi.ErrorCodeBadRequest, "invalid company id")
		return
	}

	c, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, company.ErrCompanyNotFound) {
			writeErr(w, http.StatusNotFound, oapi.ErrorCodeNotFound, "company not found")
			return
		}

		writeErr(w, http.StatusInternalServerError, oapi.ErrorCodeInternalError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, CompanyToResponse(c))
}

func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var req oapi.CreateCompanyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeErr(w, http.StatusBadRequest, oapi.ErrorCodeBadRequest, "invalid request body")
		return
	}

	params := CreateRequestToParams(req)
	c, err := h.service.CreateCompany(r.Context(), params)
	if err != nil {
		if errors.Is(err, company.ErrCompanyNameAlreadyExists) {
			writeErr(w, http.StatusConflict, oapi.ErrorCodeConflict, "company name already exists")
			return
		}

		if errors.Is(err, company.ErrInvalidCompanyNameLength) ||
			errors.Is(err, company.ErrInvalidCompanyDescriptionLength) ||
			errors.Is(err, company.ErrInvalidEmployeesCount) ||
			errors.Is(err, company.ErrInvalidCompanyType) {
			writeErr(w, http.StatusBadRequest, oapi.ErrorCodeBadRequest, err.Error())
			return
		}

		// All other errors are internal (database, kafka, etc.)
		writeErr(w, http.StatusInternalServerError, oapi.ErrorCodeInternalError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, CompanyToResponse(c))
}

func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		writeErr(w, http.StatusBadRequest, oapi.ErrorCodeBadRequest, "invalid company id")
		return
	}

	var req oapi.UpdateCompanyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeErr(w, http.StatusBadRequest, oapi.ErrorCodeBadRequest, "invalid request body")
		return
	}

	params := UpdateRequestToParams(req)
	c, err := h.service.UpdateCompany(r.Context(), id, params)
	if err != nil {
		if errors.Is(err, company.ErrCompanyNotFound) {
			writeErr(w, http.StatusNotFound, oapi.ErrorCodeNotFound, "company not found")
			return
		}

		if errors.Is(err, company.ErrCompanyNameAlreadyExists) {
			writeErr(w, http.StatusConflict, oapi.ErrorCodeConflict, "company name already exists")
			return
		}

		if errors.Is(err, company.ErrInvalidCompanyNameLength) ||
			errors.Is(err, company.ErrInvalidCompanyDescriptionLength) ||
			errors.Is(err, company.ErrInvalidEmployeesCount) ||
			errors.Is(err, company.ErrInvalidCompanyType) {
			writeErr(w, http.StatusBadRequest, oapi.ErrorCodeBadRequest, err.Error())
			return
		}

		// All other errors are internal (database, kafka, etc.)
		writeErr(w, http.StatusInternalServerError, oapi.ErrorCodeInternalError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, CompanyToResponse(c))
}

func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		writeErr(w, http.StatusBadRequest, oapi.ErrorCodeBadRequest, "invalid company id")
		return
	}

	err := h.service.DeleteCompany(r.Context(), id)
	if err != nil {
		if errors.Is(err, company.ErrCompanyNotFound) {
			writeErr(w, http.StatusNotFound, oapi.ErrorCodeNotFound, "company not found")
			return
		}

		writeErr(w, http.StatusInternalServerError, oapi.ErrorCodeInternalError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

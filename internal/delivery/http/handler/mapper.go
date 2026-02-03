package handler

import (
	"github.com/dubininme/xm-assessment/internal/domain/company"
	"github.com/dubininme/xm-assessment/pkg/gen/oapi"
)

func CompanyToResponse(c *company.Company) oapi.Company {
	oc := oapi.Company{
		Id:             c.ID(),
		Name:           c.Name().String(),
		EmployeesCount: c.EmployeesCount().Int(),
		Registered:     c.IsRegistered(),
		Type:           oapi.CompanyType(c.CompanyType().String()),
	}

	desc := c.Description().String()
	if len(desc) > 0 {
		oc.Description = &desc
	}

	return oc

}
func CreateRequestToParams(req oapi.CreateCompanyRequest) company.CreateParams {
	c := company.CreateParams{
		Name:           req.Name,
		EmployeesCount: req.EmployeesCount,
		Registered:     req.Registered,
		Type:           string(req.Type),
	}
	if req.Description != nil {
		c.Description = *req.Description
	}

	return c
}

func UpdateRequestToParams(req oapi.UpdateCompanyRequest) company.UpdateParams {
	params := company.UpdateParams{}
	if req.Name != nil {
		params.Name = req.Name
	}

	if req.EmployeesCount != nil {
		params.EmployeesCount = req.EmployeesCount
	}

	if req.Description != nil {
		params.Description = req.Description
	}

	if req.Registered != nil {
		params.Registered = req.Registered
	}

	if req.Type != nil {
		t := string(*req.Type)
		params.Type = &t
	}

	return params
}

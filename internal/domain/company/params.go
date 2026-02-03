package company

type CreateParams struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	EmployeesCount int    `json:"employees_count"`
	Registered     bool   `json:"registered"`
	Type           string `json:"type"`
}

type UpdateParams struct {
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	EmployeesCount *int    `json:"employees_count,omitempty"`
	Registered     *bool   `json:"registered,omitempty"`
	Type           *string `json:"type,omitempty"`
}

// IsEmpty checks if all fields are nil (no actual changes)
func (p UpdateParams) IsEmpty() bool {
	return p.Name == nil &&
		p.Description == nil &&
		p.EmployeesCount == nil &&
		p.Registered == nil &&
		p.Type == nil
}

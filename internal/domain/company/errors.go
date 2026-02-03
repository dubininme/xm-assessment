package company

import "errors"

var ErrInvalidCompanyNameLength = errors.New("invalid company name length")
var ErrCompanyNameAlreadyExists = errors.New("company name already exists")
var ErrInvalidCompanyDescriptionLength = errors.New("invalid company description length")
var ErrInvalidEmployeesCount = errors.New("invalid employees count")
var ErrInvalidCompanyType = errors.New("invalid company type")
var ErrCompanyNotFound = errors.New("company not found")
var ErrNoFieldsToUpdate = errors.New("at least one field must be provided for update")

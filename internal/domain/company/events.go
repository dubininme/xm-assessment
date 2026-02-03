package company

import "time"

type CompanyCreatedEvent struct {
	companyID string
	created   int64
	payload   CreateParams
}

func NewCompanyCreatedEvent(companyID string, params CreateParams) CompanyCreatedEvent {
	return CompanyCreatedEvent{
		companyID: companyID,
		created:   time.Now().Unix(),
		payload:   params,
	}
}

func (e CompanyCreatedEvent) EventName() string {
	return "CompanyCreated"
}

func (e CompanyCreatedEvent) AggregateID() string {
	return e.companyID
}

func (e CompanyCreatedEvent) CreatedAt() int64 {
	return e.created
}

func (e CompanyCreatedEvent) Payload() any {
	return struct {
		CompanyID string `json:"company_id"`
		CreateParams
	}{
		CompanyID:    e.companyID,
		CreateParams: e.payload,
	}
}

type CompanyUpdatedEvent struct {
	companyID string
	created   int64
	payload   UpdateParams
}

func NewCompanyUpdatedEvent(companyID string, params UpdateParams) CompanyUpdatedEvent {
	return CompanyUpdatedEvent{
		companyID: companyID,
		created:   time.Now().Unix(),
		payload:   params,
	}
}

func (e CompanyUpdatedEvent) EventName() string {
	return "CompanyUpdated"
}

func (e CompanyUpdatedEvent) AggregateID() string {
	return e.companyID
}

func (e CompanyUpdatedEvent) CreatedAt() int64 {
	return e.created
}

func (e CompanyUpdatedEvent) Payload() any {
	return struct {
		CompanyID string `json:"company_id"`
		UpdateParams
	}{
		CompanyID:    e.companyID,
		UpdateParams: e.payload,
	}
}

type CompanyDeletedEvent struct {
	companyID string
	created   int64
}

func NewCompanyDeletedEvent(companyID string) CompanyDeletedEvent {
	return CompanyDeletedEvent{
		companyID: companyID,
		created:   time.Now().Unix(),
	}
}

func (e CompanyDeletedEvent) EventName() string {
	return "CompanyDeleted"
}

func (e CompanyDeletedEvent) AggregateID() string {
	return e.companyID
}

func (e CompanyDeletedEvent) CreatedAt() int64 {
	return e.created
}

func (e CompanyDeletedEvent) Payload() any {
	return struct {
		CompanyID string `json:"company_id"`
	}{
		CompanyID: e.companyID,
	}
}

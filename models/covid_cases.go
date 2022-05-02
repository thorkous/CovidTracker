package models

type CovidCases struct {
	State     string  `json:"state,omitempty" validate:"required"`
	TotalCase float64 `json:"totalcase,omitempty" validate:"required"`
	Timestamp string  `json:"Timestamp,omitempty" validate:"required"`
}

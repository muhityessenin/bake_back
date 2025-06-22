package domain

import "time"

// Удаляем ReportDate struct - больше не нужен

type MarketingSource struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type SalesTeam struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type MarketingData struct {
	ID              int       `json:"id" db:"id"`
	Date            string    `json:"date" db:"date"` // Новое поле вместо ReportDateID
	SourceID        int       `json:"source_id" db:"source_id"`
	Expense         float64   `json:"expense" db:"expense"`
	Leads           int       `json:"leads" db:"leads"`
	TrialsScheduled int       `json:"trials_scheduled" db:"trials_scheduled"`
	TrialsConducted int       `json:"trials_conducted" db:"trials_conducted"`
	Payments        int       `json:"payments" db:"payments"`
	TotalAmount     float64   `json:"total_amount" db:"total_amount"`
	IsSaved         bool      `json:"is_saved" db:"is_saved"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type SalesData struct {
	ID              int       `json:"id" db:"id"`
	Date            string    `json:"date" db:"date"` // Новое поле вместо ReportDateID
	TeamID          int       `json:"team_id" db:"team_id"`
	Leads           int       `json:"leads" db:"leads"`
	TrialsScheduled int       `json:"trials_scheduled" db:"trials_scheduled"`
	TrialsConducted int       `json:"trials_conducted" db:"trials_conducted"`
	Payments        int       `json:"payments" db:"payments"`
	TotalAmount     float64   `json:"total_amount" db:"total_amount"`
	KaspiRefund     float64   `json:"kaspi_refund" db:"kaspi_refund"`
	IsSaved         bool      `json:"is_saved" db:"is_saved"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

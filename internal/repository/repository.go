package repository

import (
	"bake_backend/internal/domain"
)

type Repository interface {
	GetMarketingSources() ([]domain.MarketingSource, error)
	GetSalesTeams() ([]domain.SalesTeam, error)
	GetMarketingData(from, to string, sourceIDs []string) ([]domain.MarketingData, error)
	GetSalesData(from, to string, teamIDs []string) ([]domain.SalesData, error)
	SaveMarketingData(data *domain.MarketingData) error
	SaveSalesData(data *domain.SalesData) error
	UpdateMarketingData(data *domain.MarketingData) error
	UpdateSalesData(data *domain.SalesData) error
	GetAvailableDates() ([]string, error)
	GetAvailableMarketingDates() ([]string, error)
	GetAvailableSalesDates() ([]string, error)
}

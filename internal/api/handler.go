package api

import (
	"bake_backend/internal/domain"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// Новый метод для получения доступных дат
func (h *Handler) GetAvailableDates(w http.ResponseWriter, r *http.Request) {
	dates, err := h.repo.GetAvailableDates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dates)
}

func (h *Handler) GetMarketingSources(w http.ResponseWriter, r *http.Request) {
	sources, err := h.repo.GetMarketingSources()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sources)
}

func (h *Handler) GetSalesTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.repo.GetSalesTeams()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func (h *Handler) GetMarketingData(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	sourceIDsParam := r.URL.Query().Get("source_ids")

	var sourceIDs []string
	if sourceIDsParam != "" {
		sourceIDs = strings.Split(sourceIDsParam, ",")
	}

	data, err := h.repo.GetMarketingData(from, to, sourceIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) GetSalesData(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	teamIDsParam := r.URL.Query().Get("team_ids")

	var teamIDs []string
	if teamIDsParam != "" {
		teamIDs = strings.Split(teamIDsParam, ",")
	}

	data, err := h.repo.GetSalesData(from, to, teamIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) SaveMarketingData(w http.ResponseWriter, r *http.Request) {
	var data domain.MarketingData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.SaveMarketingData(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) SaveSalesData(w http.ResponseWriter, r *http.Request) {
	var data domain.SalesData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.SaveSalesData(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) UpdateMarketingData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var data domain.MarketingData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.ID = id
	if err := h.repo.UpdateMarketingData(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) UpdateSalesData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var data domain.SalesData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.ID = id
	if err := h.repo.UpdateSalesData(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

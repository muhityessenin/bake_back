package repository

import (
	"bake_backend/internal/domain"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetMarketingSources() ([]domain.MarketingSource, error) {
	rows, err := r.db.Query("SELECT id, name, created_at, updated_at FROM marketing_sources")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []domain.MarketingSource
	for rows.Next() {
		var s domain.MarketingSource
		if err := rows.Scan(&s.ID, &s.Name, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sources = append(sources, s)
	}
	return sources, nil
}

func (r *PostgresRepository) GetSalesTeams() ([]domain.SalesTeam, error) {
	rows, err := r.db.Query("SELECT id, name, created_at, updated_at FROM sales_teams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []domain.SalesTeam
	for rows.Next() {
		var t domain.SalesTeam
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

func (r *PostgresRepository) GetMarketingData(from, to string, sourceIDs []string) ([]domain.MarketingData, error) {
	baseQuery := "SELECT id, date, source_id, expense, leads, trials_scheduled, trials_conducted, payments, total_amount, is_saved, created_at, updated_at FROM marketing_data"
	var conditions []string
	var args []interface{}
	argIndex := 1

	if from != "" {
		conditions = append(conditions, fmt.Sprintf("date >= $%d", argIndex))
		args = append(args, from)
		argIndex++
	}

	if to != "" {
		conditions = append(conditions, fmt.Sprintf("date <= $%d", argIndex))
		args = append(args, to)
		argIndex++
	}

	if len(sourceIDs) > 0 {
		// Convert string IDs to integers and filter out invalid ones
		var validIDs []interface{}
		for _, idStr := range sourceIDs {
			if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
				validIDs = append(validIDs, id)
			}
		}

		if len(validIDs) > 0 {
			placeholders := make([]string, len(validIDs))
			for i := range validIDs {
				placeholders[i] = fmt.Sprintf("$%d", argIndex)
				args = append(args, validIDs[i])
				argIndex++
			}
			conditions = append(conditions, fmt.Sprintf("source_id IN (%s)", strings.Join(placeholders, ",")))
		}
	}

	var query string
	if len(conditions) > 0 {
		query = baseQuery + " WHERE " + strings.Join(conditions, " AND ")
	} else {
		query = baseQuery
	}
	query += " ORDER BY date DESC, source_id"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w (query: %s, args: %v)", err, query, args)
	}
	defer rows.Close()

	var data []domain.MarketingData
	for rows.Next() {
		var d domain.MarketingData
		if err := rows.Scan(&d.ID, &d.Date, &d.SourceID, &d.Expense, &d.Leads, &d.TrialsScheduled, &d.TrialsConducted, &d.Payments, &d.TotalAmount, &d.IsSaved, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	return data, nil
}

func (r *PostgresRepository) GetSalesData(from, to string, teamIDs []string) ([]domain.SalesData, error) {
	baseQuery := "SELECT id, date, team_id, leads, trials_scheduled, trials_conducted, payments, total_amount, kaspi_refund, is_saved, created_at, updated_at FROM sales_data"
	var conditions []string
	var args []interface{}
	argIndex := 1

	if from != "" {
		conditions = append(conditions, fmt.Sprintf("date >= $%d", argIndex))
		args = append(args, from)
		argIndex++
	}

	if to != "" {
		conditions = append(conditions, fmt.Sprintf("date <= $%d", argIndex))
		args = append(args, to)
		argIndex++
	}

	if len(teamIDs) > 0 {
		// Convert string IDs to integers and filter out invalid ones
		var validIDs []interface{}
		for _, idStr := range teamIDs {
			if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
				validIDs = append(validIDs, id)
			}
		}

		if len(validIDs) > 0 {
			placeholders := make([]string, len(validIDs))
			for i := range validIDs {
				placeholders[i] = fmt.Sprintf("$%d", argIndex)
				args = append(args, validIDs[i])
				argIndex++
			}
			conditions = append(conditions, fmt.Sprintf("team_id IN (%s)", strings.Join(placeholders, ",")))
		}
	}

	var query string
	if len(conditions) > 0 {
		query = baseQuery + " WHERE " + strings.Join(conditions, " AND ")
	} else {
		query = baseQuery
	}
	query += " ORDER BY date DESC, team_id"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w (query: %s, args: %v)", err, query, args)
	}
	defer rows.Close()

	var data []domain.SalesData
	for rows.Next() {
		var d domain.SalesData
		if err := rows.Scan(&d.ID, &d.Date, &d.TeamID, &d.Leads, &d.TrialsScheduled, &d.TrialsConducted, &d.Payments, &d.TotalAmount, &d.KaspiRefund, &d.IsSaved, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	return data, nil
}

func (r *PostgresRepository) SaveMarketingData(data *domain.MarketingData) error {
	if data.ID == 0 {
		return r.db.QueryRow(
			"INSERT INTO marketing_data (date, source_id, expense, leads, trials_scheduled, trials_conducted, payments, total_amount, is_saved, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
			data.Date, data.SourceID, data.Expense, data.Leads, data.TrialsScheduled, data.TrialsConducted, data.Payments, data.TotalAmount, data.IsSaved, time.Now(), time.Now(),
		).Scan(&data.ID)
	}
	_, err := r.db.Exec(
		"UPDATE marketing_data SET date=$1, source_id=$2, expense=$3, leads=$4, trials_scheduled=$5, trials_conducted=$6, payments=$7, total_amount=$8, is_saved=$9, updated_at=$10 WHERE id=$11",
		data.Date, data.SourceID, data.Expense, data.Leads, data.TrialsScheduled, data.TrialsConducted, data.Payments, data.TotalAmount, data.IsSaved, time.Now(), data.ID,
	)
	return err
}

func (r *PostgresRepository) SaveSalesData(data *domain.SalesData) error {
	if data.ID == 0 {
		return r.db.QueryRow(
			"INSERT INTO sales_data (date, team_id, leads, trials_scheduled, trials_conducted, payments, total_amount, kaspi_refund, is_saved, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
			data.Date, data.TeamID, data.Leads, data.TrialsScheduled, data.TrialsConducted, data.Payments, data.TotalAmount, data.KaspiRefund, data.IsSaved, time.Now(), time.Now(),
		).Scan(&data.ID)
	}
	_, err := r.db.Exec(
		"UPDATE sales_data SET date=$1, team_id=$2, leads=$3, trials_scheduled=$4, trials_conducted=$5, payments=$6, total_amount=$7, kaspi_refund=$8, is_saved=$9, updated_at=$10 WHERE id=$11",
		data.Date, data.TeamID, data.Leads, data.TrialsScheduled, data.TrialsConducted, data.Payments, data.TotalAmount, data.KaspiRefund, data.IsSaved, time.Now(), data.ID,
	)
	return err
}

func (r *PostgresRepository) UpdateMarketingData(data *domain.MarketingData) error {
	_, err := r.db.Exec(
		"UPDATE marketing_data SET date=$1, source_id=$2, expense=$3, leads=$4, trials_scheduled=$5, trials_conducted=$6, payments=$7, total_amount=$8, is_saved=$9, updated_at=$10 WHERE id=$11",
		data.Date, data.SourceID, data.Expense, data.Leads, data.TrialsScheduled, data.TrialsConducted, data.Payments, data.TotalAmount, data.IsSaved, time.Now(), data.ID,
	)
	return err
}

func (r *PostgresRepository) UpdateSalesData(data *domain.SalesData) error {
	_, err := r.db.Exec(
		"UPDATE sales_data SET date=$1, team_id=$2, leads=$3, trials_scheduled=$4, trials_conducted=$5, payments=$6, total_amount=$7, kaspi_refund=$8, is_saved=$9, updated_at=$10 WHERE id=$11",
		data.Date, data.TeamID, data.Leads, data.TrialsScheduled, data.TrialsConducted, data.Payments, data.TotalAmount, data.KaspiRefund, data.IsSaved, time.Now(), data.ID,
	)
	return err
}

func (r *PostgresRepository) GetAvailableMarketingDates() ([]string, error) {
	rows, err := r.db.Query("SELECT DISTINCT date FROM marketing_data ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}
	return dates, nil
}

func (r *PostgresRepository) GetAvailableSalesDates() ([]string, error) {
	rows, err := r.db.Query("SELECT DISTINCT date FROM sales_data ORDER BY date DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}
	return dates, nil
}

func (r *PostgresRepository) GetAvailableDates() ([]string, error) {
	query := `
		SELECT DISTINCT date FROM (
			SELECT date FROM marketing_data
			UNION
			SELECT date FROM sales_data
		) AS all_dates
		ORDER BY date DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}
	return dates, nil
}

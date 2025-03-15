package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type RateRepository interface {
	SaveRate(ctx context.Context, base, target string, rate float64) error
	GetLatestRates(ctx context.Context, base, target string, limit int) ([]Rate, error)
}

type Rate struct {
	BaseCurrency   string
	TargetCurrency string
	Rate           float64
	Timestamp      time.Time
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) RateRepository {
	return &postgresRepo{db: db}
}

func (p *postgresRepo) SaveRate(ctx context.Context, base, target string, rate float64) error {
	query := `
		INSERT INTO rates(base_currency, target_currency, rate)
		VALUES($1, $2, $3)`

	_, err := p.db.ExecContext(ctx, query, base, target, rate)
	if err != nil {
		return fmt.Errorf("failed to save rate: %w", err)
	}
	return nil
}

func (p *postgresRepo) GetLatestRates(ctx context.Context, base, target string, limit int) ([]Rate, error) {
	query := `
		SELECT base_currency, target_currency, rate, timestamp 
		FROM rates 
		WHERE base_currency = $1 AND target_currency = $2
		ORDER BY timestamp DESC
		LIMIT $3`

	rows, err := p.db.QueryContext(ctx, query, base, target, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query rates: %w", err)
	}
	defer rows.Close()

	var rates []Rate
	for rows.Next() {
		var r Rate
		err := rows.Scan(
			&r.BaseCurrency,
			&r.TargetCurrency,
			&r.Rate,
			&r.Timestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rate: %w", err)
		}
		rates = append(rates, r)
	}
	return rates, nil
}

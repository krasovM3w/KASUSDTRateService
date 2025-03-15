package postgres

import (
	"context"
	_ "database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepo_SaveRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepo(db)

	mock.ExpectExec("INSERT INTO rates").
		WithArgs("KAS", "USDT", 1.0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveRate(context.Background(), "KAS", "USDT", 1.0)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepo_GetLatestRates(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPostgresRepo(db)

	rows := sqlmock.NewRows([]string{"base_currency", "target_currency", "rate", "timestamp"}).
		AddRow("KAS", "USDT", 1.0, time.Now())

	mock.ExpectQuery("SELECT base_currency, target_currency, rate, timestamp FROM rates").
		WithArgs("KAS", "USDT", 1).
		WillReturnRows(rows)

	rates, err := repo.GetLatestRates(context.Background(), "KAS", "USDT", 1)
	assert.NoError(t, err)
	assert.Len(t, rates, 1)
	assert.Equal(t, "KAS", rates[0].BaseCurrency)
	assert.Equal(t, "USDT", rates[0].TargetCurrency)
	assert.Equal(t, 1.0, rates[0].Rate)
	assert.NoError(t, mock.ExpectationsWereMet())
}

package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sugyk/rest_vpn/repository"
)

var createUserQuery = `
	INSERT INTO "Users" \(telegram_id, is_admin\)
	VALUES \(\?, \?\)
	ON CONFLICT \(telegram_id\)
	DO UPDATE SET
		is_admin = EXCLUDED.is_admin;
	`

var createKeyQuery = `
	INSERT INTO "AccessKeys" \(user_id, key_value, expires_at\)
	VALUES \(\?, \?, \?\)
	`

func TestCreateUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := repository.NewRepo(sqlxDB)

	mock.ExpectExec(createUserQuery).
		WithArgs(123, true).
		WillReturnResult(sqlmock.NewResult(1, 1)) // ID, RowsAffected

	err := repo.CreateUser(repository.CreateUserParams{Telegram_id: 123, Is_admin: true})

	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestCreateKey(t *testing.T) {
	// Create mock db
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Create repo layer instance
	repo := repository.NewRepo(sqlxDB)

	// Define mock execution
	current_time := time.Now()
	mock.ExpectExec(createKeyQuery).
		WithArgs(123, "ss://asdfasdfhasdad", current_time).
		WillReturnResult(sqlmock.NewResult(1, 1)) // ID, RowsAffected

	// Try tested function
	err := repo.CreateKey(repository.CreateKeyParams{
		Telegram_id:   123,
		AccessUrl:     "ss://asdfasdfhasdad",
		Expire_months: current_time,
	},
	)

	if err != nil {
		t.Error("unexpected error", err)
	}
}

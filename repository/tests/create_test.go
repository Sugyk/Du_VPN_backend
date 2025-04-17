package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sugyk/rest_vpn/repository"
)

var query = `
	INSERT INTO "Users" \(telegram_id, is_admin\)
	VALUES \(\?, \?\)
	ON CONFLICT \(telegram_id\)
	DO UPDATE SET
		is_admin = EXCLUDED.is_admin;
	`

func TestCreateUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := repository.NewRepo(sqlxDB)

	mock.ExpectExec(query).
		WithArgs(123, true).
		WillReturnResult(sqlmock.NewResult(1, 1)) // ID, RowsAffected

	err := repo.CreateUser(repository.CreateUserParams{Telegram_id: 123, Is_admin: true})

	if err != nil {
		t.Error("unexpected error", err)
	}
}

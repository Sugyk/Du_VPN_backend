package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sugyk/rest_vpn/repository"
)

func TestCreateUser(t *testing.T) {
	// 1. Create mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewRepo(sqlxDB)

	// 2. Define test input
	params := repository.CreateUserParams{
		Telegram_id: 123456,
		Is_admin:    true,
	}

	// 3. Expect the query
	mock.ExpectExec(`INSERT INTO "Users" \(telegram_id, id_admin\) VALUES \(\$1, \$2\)`).
		WithArgs(params.Telegram_id, params.Is_admin).
		WillReturnResult(sqlmock.NewResult(1, 1)) // simulate 1 row inserted

	// 4. Call the method
	err = repo.CreateUser(params)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// 5. Ensure expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

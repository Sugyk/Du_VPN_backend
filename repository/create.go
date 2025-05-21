package repository

import "time"

type CreateUserParams struct {
	Telegram_id int  `db:"telegram_id"`
	Is_admin    bool `db:"is_admin"`
}

func (r *Repository) CreateUser(params CreateUserParams) error {
	// Create the users in db

	query := `
	INSERT INTO "Users" (telegram_id, is_admin)
	VALUES (:telegram_id, :is_admin)
	ON CONFLICT (telegram_id)
	DO UPDATE SET
		is_admin = EXCLUDED.is_admin;
	`

	_, err := r.db.NamedExec(query, params)

	if err != nil {
		return err
	}

	return nil
}

// params to repository handler CreateKey
type CreateKeyParams struct {
	TelegramId int       `db:"user_id"`
	AccessUrl  string    `db:"key_value"`
	ExpiresAt  time.Time `db:"expires_at"`
}

// Create the key for user
func (r *Repository) CreateKey(params CreateKeyParams) error {

	query := `
	INSERT INTO "AccessKeys" (user_id, key_value, expires_at)
	VALUES (:user_id, :key_value, :expires_at)
	`

	_, err := r.db.NamedExec(query, params)

	if err != nil {
		return err
	}

	return nil
}

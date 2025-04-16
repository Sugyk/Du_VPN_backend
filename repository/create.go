package repository

type CreateUserParams struct {
	Telegram_id int  `db:"telegram_id"`
	Is_admin    bool `db:"is_admin"`
}

func (r *Repository) CreateUser(params CreateUserParams) error {
	// Create the users in db

	query := `INSERT INTO "Users" (telegram_id, is_admin) VALUES (:telegram_id, :is_admin)`

	_, err := r.db.NamedExec(query, params)

	if err != nil {
		return err
	}

	return nil
}

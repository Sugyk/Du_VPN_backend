package repository

type DeleteKeyParams struct {
	KeyId int `db:"id"`
}

// delete the AccessKey from repo by id
func (r *Repository) DeleteKey(params DeleteKeyParams) error {
	query := `
	DELETE FROM "AccessKeys"
	WHERE id = :id
	`

	_, err := r.db.NamedExec(query, params)

	if err != nil {
		return err
	}

	return nil
}

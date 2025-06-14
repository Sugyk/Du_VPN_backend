package repository

import (
	"time"
)

type AccessKey struct {
	Id        int       `db:"id"`
	KeyValue  string    `db:"key_value"`
	UserId    int       `db:"user_id"`
	OutlineId int       `db:"outline_id"`
	CreatedAt time.Time `db:"created_at"`
}

type ListEntriesFilterParams struct {
	Id        int `db:"id"`
	UserId    int `db:"user_id"`
	OutlineId int `db:"outline_id"`
}

// NamedQuery using this DB. Any named placeholder parameters are replaced with fields from arg.

func (r *Repository) ListEntries(filterParams any) ([]AccessKey, error) {
	accessKeys := []AccessKey{}
	query := `
	SELECT id, key_value, user_id, outline_id, created_at
	FROM "AccessKeys"
	WHERE (:id = 0 OR id = :id)
	AND (:user_id = 0 OR user_id = :user_id)
	AND (:outline_id = 0 OR outline_id = :outline_id)
	`
	rows, err := r.db.NamedQuery(query, filterParams)
	for rows.Next() {
		var key AccessKey
		if err := rows.StructScan(&key); err != nil {
			return nil, err
		}
		accessKeys = append(accessKeys, key)
	}

	return accessKeys, err
}

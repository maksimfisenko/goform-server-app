package data

import (
	"database/sql"
	"errors"
)

type Role struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type RoleStorage struct {
	DB *sql.DB
}

func (s RoleStorage) GetByTitle(title string) (*Role, error) {
	query := `
		SELECT
    		id,
    		title
		FROM
    		users.roles
		WHERE
    		title = $1;
	`

	var role Role

	err := s.DB.QueryRow(query, title).Scan(
		&role.ID,
		&role.Title,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &role, nil
}

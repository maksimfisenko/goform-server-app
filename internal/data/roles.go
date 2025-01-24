package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, title).Scan(
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

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Storage struct {
	Roles RoleStorage
	Users UserStorage
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Roles: RoleStorage{DB: db},
		Users: UserStorage{DB: db},
	}
}

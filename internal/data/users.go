package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	ID          int64     `json:"id"`
	RoleID      int64     `json:"role_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    password  `json:"-"`
	IsActivated bool      `json:"is_activated"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int       `json:"version"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type UserStorage struct {
	DB *sql.DB
}

func (s UserStorage) Insert(user *User) error {
	query := `
		INSERT INTO
    		users.users (role_id, name, email, password_hash, is_activated)
		VALUES
    		($1, $2, $3, $4, $5) 
		RETURNING 
			id, created_at, updated_at, version
	`

	args := []interface{}{user.RoleID, user.Name, user.Email, user.Password.hash, user.IsActivated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (s UserStorage) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT
    		id, role_id, name, email, is_activated, created_at, updated_at, version
		FROM
    		users.users
		WHERE
    		id = $1;
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.IsActivated,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

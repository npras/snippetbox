package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DbPool *pgxpool.Pool
}

func (m *UserModel) Insert(name, email, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}
	qry := `
	INSERT INTO users (name, email, hashed_password, created)
	VALUES ($1,
					$2,
					$3,
					CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
	)
	RETURNING id
	`
	var id int
	err = m.DbPool.QueryRow(context.Background(), qry, name, email, string(hashedPassword)).Scan(&id)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if (pgError.Code == "23505") && strings.Contains(pgError.Message, "users_uc_email") {
				return 0, ErrDuplicateEmail
			}
		}
		return 0, err
	}
	return id, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	qry := `SELECT id, hashed_password from users where email = $1`

	err := m.DbPool.QueryRow(context.Background(), qry, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) IsExists(id int) (bool, error) {
	var exists bool
	qry := "select exists (select true from users where id = $1)"
	err := m.DbPool.QueryRow(context.Background(), qry, id).Scan(&exists)
	return exists, err
}

package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

//

func (m *SnippetModel) Insert(title string, content string, expiresAt int) (int, error) {
	qry := `
	INSERT INTO snippets (title, content, created_at, expires_at)
	VALUES ($1,
					$2,
					CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
					(CURRENT_TIMESTAMP AT TIME ZONE 'UTC') + (INTERVAL '1 day' * $3)
	)
	RETURNING id
	`
	var id int
	err := m.DB.QueryRow(qry, title, content, expiresAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//

func (m *SnippetModel) Get(id int) (Snippet, error) {
	var s Snippet
	stmt := `
	SELECT id, title, content, created_at, expires_at
		FROM snippets
		WHERE
			expires_at > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
			AND id = $1
	`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		}
	}
	return s, nil
}

//

func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}

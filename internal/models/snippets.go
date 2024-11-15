package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Snippet struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SnippetModel struct {
	DbPool *pgxpool.Pool
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
	err := m.DbPool.QueryRow(context.Background(), qry, title, content, expiresAt).Scan(&id)
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
			id = $1
			AND expires_at > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
	`
	err := m.DbPool.QueryRow(context.Background(), stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		}
	}
	return s, nil
}

//

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `
	SELECT id, title, content, created_at, expires_at
		FROM snippets
		WHERE
			expires_at > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		ORDER BY id DESC
		LIMIT 10	
	`
	rows, err := m.DbPool.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var snippets []Snippet

	for rows.Next() {
		var s Snippet
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.ExpiresAt)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

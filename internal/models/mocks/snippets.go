package mocks

import (
	"time"

	"github.com/npras/snippetbox/internal/models"
)

var mockSnippet = models.Snippet{
	ID:        1,
	Title:     "some titl",
	Content:   "some contt",
	CreatedAt: time.Now(),
	ExpiresAt: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expiresAt int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}

package mocks

import "github.com/npras/snippetbox/internal/models"

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) (int, error) {
	switch email {
	case "dupe@example.com":
		return 0, models.ErrDuplicateEmail
	default:
		return 2, nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) IsExists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

package user_repository

import (
	"database/sql"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (us *UserStore) GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, username, password FROM users WHERE username = $1"
	row := us.db.QueryRow(query, username)

	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

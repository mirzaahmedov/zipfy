package postgres

import (
	"database/sql"

	"zipfy/internal/model"
	"zipfy/internal/store"
)

type UserStore struct {
	db *sql.DB
}

func (ps *PostgresStore) User() store.UserStore {
	return &UserStore{ps.db}
}

func (us *UserStore) Create(user *model.User) (*model.User, error) {
	err := us.db.
		QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email, password", user.Name, user.Email, user.Password).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (us *UserStore) GetByID(id string) (*model.User, error) {
	user := new(model.User)
	err := us.db.
		QueryRow("SELECT id, name, email, password FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (us *UserStore) GetByEmail(email string) (*model.User, error) {
	user := new(model.User)
	err := us.db.
		QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

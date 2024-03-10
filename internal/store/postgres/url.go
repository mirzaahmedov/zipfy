package postgres

import (
	"database/sql"
	"errors"

	"zipfy/internal/generate"
	"zipfy/internal/model"
	"zipfy/internal/store"
)

type URLStore struct {
	db *sql.DB
}

func (ps *PostgresStore) URL() store.URLStore {
	return &URLStore{
		db: ps.db,
	}
}

func (ur *URLStore) Create(url *model.URL) (*model.URL, error) {
	err := ur.db.
		QueryRow("INSERT INTO urls (url, user_id) VALUES ($1, $2) RETURNING id, url, user_id, created", url.URL, url.UserID).
		Scan(&url.ID, &url.URL, &url.UserID, &url.Created)
	if err != nil {
		return nil, err
	}
	url.Short = generate.ShortenURL(url.ID)
	return url, nil
}

func (ur *URLStore) GetByID(id int) (*model.URL, error) {
	url := new(model.URL)
	err := ur.db.
		QueryRow("SELECT id, url, user_id, created FROM urls WHERE id = $1", id).
		Scan(&url.ID, &url.URL, &url.UserID, &url.Created)
	if err != nil {
		return nil, err
	}
	url.Short = generate.ShortenURL(url.ID)
	return url, nil
}

func (ur *URLStore) GetByUserID(id string) ([]*model.URL, error) {
	userURLs := make([]*model.URL, 0)
	rows, err := ur.db.Query("SELECT id, url, user_id, created FROM urls WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		url := new(model.URL)
		err := rows.Scan(&url.ID, &url.URL, &url.UserID, &url.Created)
		if err != nil {
			return nil, err
		}
		url.Short = generate.ShortenURL(url.ID)
		userURLs = append(userURLs, url)
	}
	return userURLs, nil
}

func (ur *URLStore) Update(id int, newURL string) (*model.URL, error) {
	url := new(model.URL)

	err := ur.db.QueryRow("UPDATE urls SET url = $1 WHERE id = $2 RETURNING id, url, user_id, created", newURL, id).
		Scan(&url.ID, &url.URL, &url.UserID, &url.Created)
	if err != nil {
		return nil, err
	}

	url.Short = generate.ShortenURL(url.ID)
	return url, nil
}

func (ur *URLStore) Delete(id int) error {
	result, err := ur.db.Exec("DELETE FROM urls WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("No rows affected")
	}
	return nil
}

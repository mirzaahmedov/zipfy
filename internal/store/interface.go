package store

import (
	"zipfy/internal/model"
)

type Store interface {
	User() UserStore
	URL() URLStore
	Open() error
	Close() error
}

type UserStore interface {
	Create(user *model.User) (*model.User, error)
	GetByID(id string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}
type URLStore interface {
	Create(url *model.URL) (*model.URL, error)
	GetByID(id int) (*model.URL, error)
	GetByUserID(id string) ([]*model.URL, error)
	Update(id int, newURL string) (*model.URL, error)
	Delete(id int) error
}

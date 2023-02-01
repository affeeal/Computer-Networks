package sqlstore

import (
	"database/sql"

	"github.com/affeeal/lab9/internal/app/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db             *sql.DB
	UserRepository *UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.UserRepository == nil {
		s.UserRepository = &UserRepository{
			store: s,
		}
	}

	return s.UserRepository
}

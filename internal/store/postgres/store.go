package postgres

import (
	"database/sql"
	"fmt"

	"zipfy/internal/store"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db         *sql.DB
	connString string
}

func NewPostgresStore(user, pwd, host, port, dbname, sslmode string) store.Store {
	return &PostgresStore{
		connString: "postgresql://" + user + ":" + pwd + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode,
	}
}

func (ps *PostgresStore) Open() error {
	db, err := sql.Open("postgres", ps.connString)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("Error pinging database: %v", err)
	}
	ps.db = db
	return nil
}
func (ps *PostgresStore) Close() error {
	return ps.db.Close()
}

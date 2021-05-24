package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"database/sql"

	// Postgres db library loading
	_ "github.com/lib/pq"
)

// PostgresqlStore is a storage engine that writes to postgres
type PostgresqlStore struct {
	db *sql.DB
}

// NewPostgresqlClient creates a new db client object
func NewPostgresqlClient(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(`
		CREATE TABLE IF NOT EXISTS users (
			id varchar(255) NOT NULL,
			username varchar(255) NOT NULL,
			access varchar(255) NOT NULL,
			refresh varchar(255) NOT NULL,
			updated timestamp with time zone NOT NULL,
			PRIMARY KEY(id)
		)
	`)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	return db
}

// NewPostgresqlStore creates new store
func NewPostgresqlStore(db *sql.DB) PostgresqlStore {
	return PostgresqlStore{
		db: db,
	}
}

// Ping will check if the connection works right
func (s PostgresqlStore) Ping(ctx context.Context) error {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.PingContext(ctx)
}

// WriteUser will write a user object to postgres
func (s PostgresqlStore) WriteUser(user User) {
	_, err := s.db.Exec(
		`
			INSERT INTO users
				(id, username, access, refresh, updated)
				VALUES($1, $2, $3, $4, $5)
			ON CONFLICT(id)
			DO UPDATE set username=EXCLUDED.username, access=EXCLUDED.access, refresh=EXCLUDED.refresh, updated=EXCLUDED.updated
		`,
		user.ID,
		user.Username,
		user.AccessToken,
		user.RefreshToken,
		user.Updated,
	)
	if err != nil {
		panic(err)
	}
}

// GetUser will load a user from postgres
func (s PostgresqlStore) GetUser(id string) *User {
	var username string
	var access string
	var refresh string
	var updated time.Time

	err := s.db.QueryRow(
		"SELECT username, access, refresh, updated FROM users WHERE id=$1",
		id,
	).Scan(
		&username,
		&access,
		&refresh,
		&updated,
	)
	switch {
	case err == sql.ErrNoRows:
		panic(fmt.Errorf("no user with id %s", id))
	case err != nil:
		panic(fmt.Errorf("query error: %v", err))
	}
	user := User{
		ID:           id,
		Username:     strings.ToLower(username),
		AccessToken:  access,
		RefreshToken: refresh,
		Updated:      updated,
		store:        s,
	}

	return &user
}

// TODO: Not Implemented
func (s PostgresqlStore) DeleteUser(id string) bool {
	return true
}

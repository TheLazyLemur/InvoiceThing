package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type liteDB struct {
	db *sql.DB
}

func NewLiteDB() IDB {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	return &liteDB{
		db: db,
	}
}

func (ldb *liteDB) GetUser(ctx context.Context, email string) (string, string, error) {
	rows, err := ldb.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var em, pwrd string

		if err := rows.Scan(&id, &em, &pwrd); err != nil {
			return "", "", err
		}

		return email, pwrd, nil
	}

	return "", "", ErrUserNotFound
}

func (ldb *liteDB) CreateUser(ctx context.Context, email string, password string) error {
	_, err := ldb.db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, password)
	if err != nil {
		return err
	}

	return nil
}

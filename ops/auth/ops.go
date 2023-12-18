package auth

import (
	"context"
	"invoicething/database"
)

func Login(ctx context.Context, db database.IDB, email string, password string) (string, error) {
	_, pwrd, err := db.GetUser(ctx, email)
	if err != nil {
		return "", err
	}

	if pwrd == password {
		return email, nil
	} else {
		return "", database.ErrWrongPassword
	}
}

func CreateUser(ctx context.Context, db database.IDB, email string, password string) error {
	_, _, err := db.GetUser(ctx, email)
	if err == nil {
		return database.ErrUserExists
	}

	return db.CreateUser(ctx, email, password)
}

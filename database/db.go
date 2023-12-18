package database

import (
	"context"
	"errors"
)

var (
	ErrUserExists    = errors.New("User already exists")
	ErrUserNotFound  = errors.New("User not found")
	ErrWrongPassword = errors.New("Wrong password")
)

type IDB interface {
	GetUser(ctx context.Context, email string) (string, string, error)
	CreateUser(ctx context.Context, email, password string) error
}

package user

import (
	"context"
	"database/sql"
)

type UserUC struct {
	db *sql.DB
}

func NewUserUC(db *sql.DB) *UserUC {
	return &UserUC{db}
}

func (u *UserUC) Create(ctx context.Context, req User) error {
	return u.create(ctx, req)
}

func (u *UserUC) List(ctx context.Context) ([]User, error) {
	resp, err := u.list(ctx)
	if err != nil {
		return []User{}, err
	}
	return resp, err
}

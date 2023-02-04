package user

import (
	"context"
	"database/sql"
)

type userUC struct {
	db *sql.DB
}

func NewUserUC(db *sql.DB) *userUC {
	return &userUC{db}
}

func (u *userUC) Create(ctx context.Context, req User) error {
	return u.create(ctx, req)
}

func (u *userUC) List(ctx context.Context) ([]User, error) {
	resp, err := u.list(ctx)
	if err != nil {
		return []User{}, err
	}
	return resp, err
}

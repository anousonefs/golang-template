package user

import (
	"context"
	"database/sql"
)

type UserUsecase struct {
	db *sql.DB
}

func NewUserUsecase(db *sql.DB) *UserUsecase {
	return &UserUsecase{db}
}

func (u *UserUsecase) Create(ctx context.Context, req User) error {
	return u.create(ctx, req)
}

func (u *UserUsecase) List(ctx context.Context) ([]User, error) {
	resp, err := u.list(ctx)
	if err != nil {
		return []User{}, err
	}
	return resp, err
}

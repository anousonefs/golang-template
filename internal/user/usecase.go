package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/anousoneFS/clean-architecture/helper"
)

type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) UserUsecase {
	return UserUsecase{repo}
}

func (u *UserUsecase) Create(ctx context.Context, req User) error {
	return u.repo.Create(ctx, req)
}

func (u *UserUsecase) List(ctx context.Context) ([]User, error) {
	res, err := u.repo.List(ctx)
	if err != nil {
		return []User{}, err
	}
	return res, err
}

func (u *UserUsecase) Get(ctx context.Context, id string) (res *User, err error) {
	res, err = u.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.AppError{Code: 406, Message: "wowow"}
			// return nil, helper.ErrUnauthorize
		}
		return nil, err
	}
	return res, err
}

func (u *UserUsecase) GetByUsername(ctx context.Context, username string) (res *User, err error) {
	res, err = u.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, helper.AppError{Code: 401}
		}
		return nil, err
	}
	return res, err
}

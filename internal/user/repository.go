package user

import (
	"context"
	"database/sql"

	"github.com/anousoneFS/clean-architecture/config"
)

type UserRepo struct {
	db *sql.DB
}

func New(db *sql.DB) UserRepo {
	return UserRepo{db: db}
}

func (r UserRepo) List(ctx context.Context) ([]User, error) {
	query, args, err := config.Psql().
		Select("name", "age", "phone").
		From("users").
		ToSql()
	if err != nil {
		return []User{}, err
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	defer rows.Close()
	resp := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.Name, &i.Age, &i.Phone); err != nil {
			return []User{}, err
		}
		resp = append(resp, i)
	}
	if err != nil {
		return []User{}, err
	}
	return resp, nil
}

func (r UserRepo) Create(ctx context.Context, req User) error {
	query, args, err := config.Psql().
		Insert("users").
		Columns("name", "age").
		Values(req.Name, req.Age).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepo) Get(ctx context.Context, id string) (res *User, err error) {
	query, args, err := config.Psql().
		Select("name", "age", "phone").
		From("users").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, err
	}
	var i User
	row := r.db.QueryRowContext(ctx, query, args...)
	if err := row.Scan(&i.Name, &i.Age, &i.Phone); err != nil {
		return nil, err
	}
	return &i, nil
}

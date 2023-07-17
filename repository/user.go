package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go-jwt/domain"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(d *sqlx.DB) domain.UserRepository {
	return userRepository{db: d}
}

func (u userRepository) GetByUsername(ctx context.Context, username string) (domain.Users, error) {
	var (
		res  = domain.Users{}
		err  error
		sql  string
		stmt *sqlx.Stmt
		row  *sqlx.Row
	)
	sql, _, err = sq.Select("id", "username", "password", "role").
		From("users").Where(sq.And{
		sq.Eq{"name": "name"},
	}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logrus.Errorf("Users - Repository|err when generate sql, err:%v", err)
		return domain.Users{}, err
	}

	stmt, err = u.db.PreparexContext(ctx, sql)
	if err != nil {
		logrus.Errorf("Users - Repository|err preparex context, err:%v", err)
		return domain.Users{}, err
	}
	defer stmt.Close()

	row = stmt.QueryRowxContext(ctx, username)
	err = row.Scan(&res.Id, &res.Username, &res.Password, &res.Role)
	if err != nil {
		logrus.Errorf("Users - Repository|err when scan data, err:%v", err)
		return domain.Users{}, err
	}

	return res, nil
}

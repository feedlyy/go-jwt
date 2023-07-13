package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go-jwt/domain"
)

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(d *sqlx.DB) domain.ProductRepository {
	return productRepository{db: d}
}

func (p productRepository) GetByName(ctx context.Context, name string) (domain.Products, error) {
	var (
		res  = domain.Products{}
		err  error
		sql  string
		stmt *sqlx.Stmt
		row  *sqlx.Row
	)
	sql, _, err = sq.Select("id", "name", "qty", "description", "created_at", "updated_at").
		From("products").Where(sq.And{
		sq.Eq{"name": "name"},
	}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logrus.Errorf("Products - Repository|err when generate sql, err:%v", err)
		return domain.Products{}, err
	}

	stmt, err = p.db.PreparexContext(ctx, sql)
	if err != nil {
		logrus.Errorf("Products - Repository|err preparex context, err:%v", err)
		return domain.Products{}, err
	}
	defer stmt.Close()

	row = stmt.QueryRowxContext(ctx, name)
	err = row.Scan(&res.Id, &res.Name, &res.Qty, &res.Description, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		logrus.Errorf("Products - Repository|err when scan data, err:%v", err)
		return domain.Products{}, err
	}

	return res, nil
}

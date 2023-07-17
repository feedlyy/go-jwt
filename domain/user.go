package domain

import (
	"context"
	"database/sql"
	"github.com/golang-jwt/jwt/v5"
)

type Users struct {
	Id        int            `json:"id"`
	Username  sql.NullString `json:"username"`
	Name      sql.NullString `json:"name"`
	Password  sql.NullString `json:"password"`
	Role      sql.NullString `json:"role"`
	CreatedAt sql.NullTime   `json:"createdAt"`
	UpdatedAt sql.NullTime   `json:"updatedAt"`
}

type JwtClaims struct {
	jwt.RegisteredClaims
	Username string `json:"Username"`
	Role     string `json:"Role"`
}

type JwtCredentials struct {
	ApplicationName string
	SecretKey       string
	Expired         int
}

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (Users, error)
}

type UserService interface {
	Authentication(ctx context.Context, username, password string) (string, error)
}

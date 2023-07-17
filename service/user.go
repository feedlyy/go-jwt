package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"go-jwt/domain"
	"go-jwt/helpers"
	"time"
)

type userService struct {
	userRepo domain.UserRepository
	jwt      domain.JwtCredentials
}

func NewUserService(u domain.UserRepository, j domain.JwtCredentials) domain.UserService {
	return userService{userRepo: u, jwt: j}
}

func (u userService) Authentication(ctx context.Context, username, password string) (string, error) {
	var (
		user        = domain.Users{}
		err         error
		claims      = domain.JwtClaims{}
		token       *jwt.Token
		signedToken string
	)

	user, err = u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	err = helpers.ValidatePassword(user.Password.String, password)
	if err != nil {
		logrus.Errorf("Users - Service|Password err, err:%v", err)
		return "", errors.New(helpers.IncorrectPassword)
	}
	fmt.Println(user.Username.String)

	claims = domain.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    u.jwt.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(u.jwt.Expired) * time.Hour)),
		},
		Username: user.Username.String,
		Role:     user.Role.String,
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(u.jwt.SecretKey))
	if err != nil {
		logrus.Errorf("Users - Service|err generate signedToken, err:%v", err)
		return "", err
	}

	return signedToken, nil
}

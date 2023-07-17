package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"go-jwt/domain"
	"net/http"
	"strings"
)

type Middleware struct {
	jwt domain.JwtCredentials
}

func NewMiddleware(j domain.JwtCredentials) Middleware {
	return Middleware{
		jwt: j,
	}
}

func (m Middleware) AuthJWTMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		ctx := context.Background()

		authHeader := request.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		writer.Header().Set("Content-Type", "application/json")

		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logrus.Errorf("Middleware|Check signin method")
				return nil, fmt.Errorf("Signing method invalid")
			} else if method != jwt.SigningMethodHS256 {
				logrus.Errorf("Middleware|Check signin method")
				return nil, fmt.Errorf("Signing method invalid")
			}

			return []byte(m.jwt.SecretKey), nil
		})
		if err != nil {
			switch {
			case strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()):
				http.Error(writer, err.Error(), http.StatusUnauthorized)
				return
			default:
				logrus.Errorf("Middleware|Err jwt parse, err:%v", err)
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// token.Claims are refering to data from what I set on user service: jwt.NewWithClaims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		// Add the token to the request context
		ctx = context.WithValue(ctx, "userInfo", claims)

		// If authentication succeeded, call the next handler with the modified context
		next(writer, request.WithContext(ctx), params)
	}
}

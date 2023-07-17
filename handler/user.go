package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-jwt/domain"
	"go-jwt/helpers"
	"net/http"
	"time"
)

type userHandler struct {
	userService domain.UserService
	timeout     time.Duration
}

func NewUserHandler(u domain.UserService, t time.Duration) userHandler {
	return userHandler{
		userService: u,
		timeout:     t,
	}
}

func (u userHandler) Authentication(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err  error
		resp = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
		token    string
		username = r.PostFormValue("username")
		password = r.PostFormValue("password")
	)
	w.Header().Set("Content-Type", "application/json")

	switch {
	case username == "":
		resp.Status = helpers.FailMsg
		resp.Message = fmt.Sprintf(helpers.MissingBody, "username")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	case password == "":
		resp.Status = helpers.FailMsg
		resp.Message = fmt.Sprintf(helpers.MissingBody, "password")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	token, err = u.userService.Authentication(ctx, username, password)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Data = err.Error()

		switch {
		case err == sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resp)
			return
		case err.Error() == helpers.IncorrectPassword:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		default:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	resp.Data = map[string]interface{}{
		"token": token,
	}
	json.NewEncoder(w).Encode(resp)
	return
}

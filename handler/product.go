package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go-jwt/domain"
	"go-jwt/helpers"
	"net/http"
	"time"
)

type productHandler struct {
	productService domain.ProductService
	timeout        time.Duration
}

func NewProductHandler(p domain.ProductService, t time.Duration) productHandler {
	handler := &productHandler{
		productService: p,
		timeout:        t,
	}

	return *handler
}

func (p productHandler) GetProductByName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err  error
		resp = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
		res  = domain.Products{}
		name = r.URL.Query().Get("name")
		ctxx = context.Background()
		// userInfo = r.Context().Value("userInfo").(jwt.MapClaims) -- for get user logged info
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(ctxx, p.timeout)
	defer cancel()

	res, err = p.productService.GetProductByName(ctx, name)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Data = err.Error()

		switch {
		case err == sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resp)
			return
		default:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	resp.Data = res
	json.NewEncoder(w).Encode(resp)
	return
}

package transport

import (
	"net/http"
	"context"
	"fmt"
	"encoding/json"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"


	localEndpoint "local.com/13sai/go-kit-demo/endpoint"
)

type errorWrapper struct {
	Error string `json:"errors"`
}

func decodeLoginRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var login localEndpoint.LoginRequest
	err := json.NewDecoder(req.Body).Decode(&login)
	fmt.Println(login)
	if err != nil {
		return nil, err
	}
	return &login, nil
}

func decodeRegisterRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var register localEndpoint.RegisterRequest
	err := json.NewDecoder(req.Body).Decode(&register)
	if err != nil {
		return nil, err
	}
	return register, nil
}

func encodeHTTPResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(res)
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	fmt.Println("errorEncoder", err.Error())
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func NewHttpHandle(ctx context.Context, end *localEndpoint.UserEndpoints) http.Handler {
	ops := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
	}

	m := mux.NewRouter()
	m.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		end.RegisterEndpoint,
		decodeRegisterRequest,
		encodeHTTPResponse,
		ops...,
	))
	m.Methods("POST").Path("/login").Handler(kithttp.NewServer(
		end.LoginEndPoint,
		decodeLoginRequest,
		encodeHTTPResponse,
		ops...,
	))


	return m
}
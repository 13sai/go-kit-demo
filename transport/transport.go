package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"go-kit-demo/endpoint"
	localEndpoint "go-kit-demo/endpoint"
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

func NewHttpHandle(ctx context.Context, end *localEndpoint.RegisterEndpoints) http.Handler {
	ops := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
	}

	m := mux.NewRouter()

	m.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		end.HealthCheckEndpoint,
		decodeRegisterRequest,
		encodeHTTPResponse,
		ops...,
	))

	m.Methods("GET").Path("/discovery/name").Handler(kithttp.NewServer(
		end.DiscoveryEndpoint,
		decodeDiscoveryRequest,
		encodeHTTPResponse,
		ops...,
	))
	return m
}

func decodeDiscoveryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	serviceName := r.URL.Query().Get("serviceName")

	if serviceName == "" {
		return nil, errors.New("invalid request parameter")
	}
	return endpoint.DiscoveryRequest{
		ServiceName: serviceName,
	}, nil
}

func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpoint.HealthRequest{}, nil
}

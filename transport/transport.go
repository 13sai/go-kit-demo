package transport

import (
	"net/http"
	"context"
	"strconv"
	"fmt"
	"reflect"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"


	"local.com/13sai/go-kit-demo/service"
	localEndpoint "local.com/13sai/go-kit-demo/endpoint"
)

type errorWrapper struct {
	Error string `json:"errors"`
}

func decodeHTTPLoginRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var login service.Login
	err := json.NewDecoder(req.Body).Decode(&login)
	if err != nil {
		return nil, err
	}
	return login, nil
}

func decodeHTTPAddRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var (
		in service.Add
		err error
	)

	in.A, err = strconv.Atoi(req.FormValue("a"))

	if err != nil {
		return in, err
	}
	in.B, err = strconv.Atoi(req.FormValue("b"))
	if err != nil {
		return in, err
	}
	return in, nil
}

func encodeHTTPResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	if f, ok := res.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("uuid", ctx.Value(service.ContextReqUUid).(string))
	return json.NewEncoder(w).Encode(res)
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	fmt.Println("errorEncoder", err.Error())
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func NewHttpHandle(end localEndpoint.EndPointServer) http.Handler {
	ops := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
		kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			UUID := (uuid.NewV4()).String()
			fmt.Println(UUID, reflect.TypeOf(UUID))
			ctx = context.WithValue(ctx, service.ContextReqUUid, UUID)
			ctx = context.WithValue(ctx, viper.GetString("jwt.name"), request.Header.Get("Authorization"))
			return ctx
		}),
	}

	m := http.NewServeMux()
	m.Handle("/sum", kithttp.NewServer(
		end.AddEndPoint,
		decodeHTTPAddRequest,
		encodeHTTPResponse,
		ops...,
	))
	m.Handle("/login", kithttp.NewServer(
		end.LoginEndPoint,
		decodeHTTPLoginRequest,
		encodeHTTPResponse,
		ops...,
	))


	return m
}
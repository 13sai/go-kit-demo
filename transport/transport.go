package transport

import (
	"net/http"
	"context"
	"strconv"
	"fmt"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"local.com/13sai/go-kit-demo/service"
	localEndpoint "local.com/13sai/go-kit-demo/endpoint"
)

type errorWrapper struct {
	Error string `json:"errors"`
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

func encodeHTTPAddResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	if f, ok := res.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	}

	m := http.NewServeMux()
	m.Handle("/sum", kithttp.NewServer(
		end.AddEndPoint,
		decodeHTTPAddRequest,
		encodeHTTPAddResponse,
		ops...,
	))

	return m
}
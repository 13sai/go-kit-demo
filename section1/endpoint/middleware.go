package endpoint

import (
	"context"
	"fmt"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/spf13/viper"

	"go-kit-demo/util"
)

func AuthMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint{
		return func (ctx context.Context, request interface{})  (res interface{}, err error){
			token := fmt.Sprint(ctx.Value(viper.GetString("jwt.name")))
			if token == "" {
				err = errors.New("please login")
				return
			}

			jwt, err := util.ParseToken(token)
			if err != nil {
				return
			}

			fmt.Println(jwt)

			if v, ok := jwt["Name"]; ok {
				ctx = context.WithValue(ctx, "name", v)
			}

			return next(ctx, request)
		}
	}
}
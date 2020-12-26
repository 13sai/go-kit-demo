package endpoint

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

var ErrLimitExceed = errors.New("Rate limit exceed!")

func NewTokenBucketLimit(l *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (res interface{}, err error) {
			if !l.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

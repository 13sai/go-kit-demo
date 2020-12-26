package transport

import (
	"context"
	"goKitRPC/endpoint"
	"goKitRPC/grpc/pb"

	"github.com/go-kit/kit/transport/grpc"
)

type gRpcServer struct {
	hand grpc.Handler
}

func (g *gRpcServer) CheckPassword(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, res, err := g.hand.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return res.(*pb.LoginResponse), nil
}

func NewUserServer(ctx context.Context, endpoint endpoint.UEndpoint) pb.UserServiceServer {
	return &gRpcServer{
		hand: grpc.NewServer(
			endpoint.UserEndpoint,
			DecodeLoginRequest,
			EncodeLoginRes,
		),
	}
}

func DecodeLoginRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.LoginRequest)
	return endpoint.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}, nil
}

func EncodeLoginRes(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.LoginRes)
	resStr := ""
	if resp.Ret {
		resStr = "success"
	}

	errStr := ""
	if resp.Err != nil {
		errStr = resp.Err.Error()
	}

	return &pb.LoginResponse{Ret: resStr, Err: errStr}, nil
}

package user

import (
	"context"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	gt "github.com/go-kit/kit/transport/grpc"
	pb "gokitgrpc/user/pb"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/user").Handler(httptransport.NewServer(
		endpoints.CreateUser,
		decodeUserReq,
		encodeResponse,
	))

	r.Methods("POST").Path("/auth").Handler(httptransport.NewServer(
		endpoints.AuthenticateUser,
		decodeAuthenticateUserReq,
		encodeResponse,
	))

	r.Methods("GET").Path("/user").Handler(httptransport.NewServer(
		endpoints.ListUser,
		decodeListUserReq,
		encodeResponse,
	))

	r.Methods("GET").Path("/user/{id}").Handler(httptransport.NewServer(
		endpoints.ViewUser,
		decodeViewUserReq,
		encodeResponse,
	))

	r.Methods("PUT").Path("/user/{id}").Handler(httptransport.NewServer(
		endpoints.UpdateUser,
		decodeUpdateUserReq,
		encodeResponse,
	))

	r.Methods("DELETE").Path("/user/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteUser,
		decodeDeleteUserReq,
		encodeResponse,
	))

	return r

}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type gRPCServer struct {
	view_user gt.Handler
	list_user gt.Handler
	create_user gt.Handler
	update_user gt.Handler
	delete_user gt.Handler
	auth gt.Handler
}

// NewGRPCServer initializes a new gRPC server
func NewGRPCServer(endpoints Endpoints, logger log.Logger) pb.UserGrpcServer {
    return &gRPCServer{
		view_user: gt.NewServer(
            endpoints.GViewUser,
            decodeGViewUserReq,
            encodeGUser,
        ),
		create_user: gt.NewServer(
            endpoints.GCreateUser,
            decodeGCreateUserReq,
            encodeGResponse,
        ),
    }
}

func (s *gRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.Response, error) {
    _, resp, err := s.create_user.ServeGRPC(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp.(*pb.Response), nil
}

func (s *gRPCServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Response, error) {
    _, resp, err := s.delete_user.ServeGRPC(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp.(*pb.Response), nil
}

func (s *gRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.Response, error) {
    _, resp, err := s.update_user.ServeGRPC(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp.(*pb.Response), nil
}

func (s *gRPCServer) ViewUser(ctx context.Context, req *pb.ViewUserRequest) (*pb.UserResponse, error) {
    _, resp, err := s.view_user.ServeGRPC(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp.(*pb.UserResponse), nil
}

func (s *gRPCServer) ListUsers(ctx context.Context, req *pb.ListUserRequest)  (*pb.UsersResponse, error) {
	_, resp, err := s.list_user.ServeGRPC(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp.(*pb.UsersResponse), nil
}

func (s *gRPCServer) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthenticateResponse, error) {
    _, resp, err := s.auth.ServeGRPC(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp.(*pb.AuthenticateResponse), nil
}
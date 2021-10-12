package user

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"fmt"
	pb "gokitgrpc/user/pb"
)

type Endpoints struct {
	CreateUser			endpoint.Endpoint
	ViewUser			endpoint.Endpoint
	UpdateUser			endpoint.Endpoint
	DeleteUser			endpoint.Endpoint
	ListUser			endpoint.Endpoint
	AuthenticateUser	endpoint.Endpoint

	GCreateUser			endpoint.Endpoint
	GViewUser			endpoint.Endpoint
	GUpdateUser			endpoint.Endpoint
	GDeleteUser			endpoint.Endpoint
	GListUser			endpoint.Endpoint
	GAuthenticateUser	endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateUser: makeCreateUserEndpoint(s),
		ViewUser: makeViewUserEndpoint(s),
		UpdateUser: makeUpdateUserEndpoint(s),
		DeleteUser: makeDeleteUserEndpoint(s),
		ListUser: makeListUserEndpoint(s),
		AuthenticateUser: makeAuthenticateUserEndpoint(s),

		GCreateUser: makeGCreateUserEndpoint(s),
		GViewUser: makeGViewUserEndpoint(s),
		GUpdateUser: makeGUpdateUserEndpoint(s),
		GDeleteUser: makeGDeleteUserEndpoint(s),
		GListUser: makeGListUserEndpoint(s),
		GAuthenticateUser: makeGAuthenticateUserEndpoint(s),
	}
}

func makeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		ok, err := s.CreateUser(ctx, req.Name, req.Phone, req.Email, req.Password)
		return CreateUserResponse{Ok: ok}, err
	}
}

func makeViewUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ViewUserRequest)
		user, err := s.ViewUser(ctx, req.Id)
		
		var jsonData User	
		json.Unmarshal([]byte(user), &jsonData)
		return jsonData, err
	}
}

func makeUpdateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUserRequest)
		ok, err := s.UpdateUser(ctx, req.Id, req.Name, req.Phone, req.Email)
		return UpdateUserResponse{Ok: ok}, err
	}
}

func makeDeleteUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteUserRequest)
		ok, err := s.DeleteUser(ctx, req.Id)
		return DeleteUserResponse{Ok: ok}, err
	}
}

func makeListUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListUserRequest)
		users, err := s.ListUser(ctx, req.Limit, req.Offset)
		
		var jsonData []*User	
		json.Unmarshal([]byte(users), &jsonData)
		return jsonData, err
	}
}

func makeAuthenticateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthenticateUserRequest)
		access_token, err := s.AuthenticateUser(ctx, req.Phone, req.Password)
		fmt.Println("Auth User")
		return AuthenticateUserResponse{
			AccessToken: access_token,
		}, err
	}
}


// GRPC
func makeGCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pb.CreateUserRequest)
		ok, err := s.CreateUser(ctx, req.Name, req.Phone, req.Email, req.Password)
		return CreateUserResponse{Ok: ok}, err
	}
}

func makeGViewUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pb.ViewUserRequest)
		user, err := s.ViewUser(ctx, req.Id)
		
		var jsonData User	
		json.Unmarshal([]byte(user), &jsonData)
		return jsonData, err
	}
}

func makeGUpdateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pb.UpdateUserRequest)
		ok, err := s.UpdateUser(ctx, req.Id, req.Name, req.Phone, req.Email)
		return UpdateUserResponse{Ok: ok}, err
	}
}

func makeGDeleteUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pb.DeleteUserRequest)
		ok, err := s.DeleteUser(ctx, req.Id)
		return DeleteUserResponse{Ok: ok}, err
	}
}

func makeGListUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pb.ListUserRequest)
		users, err := s.ListUser(ctx, req.Limit, req.Offset)
		
		var jsonData []*User	
		json.Unmarshal([]byte(users), &jsonData)
		return jsonData, err
	}
}

func makeGAuthenticateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pb.AuthenticateUserRequest)
		access_token, err := s.AuthenticateUser(ctx, req.Phone, req.Password)
		fmt.Println("Auth User")
		return AuthenticateUserResponse{
			AccessToken: access_token,
		}, err
	}
}
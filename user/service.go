package user

import "context"

type Service interface {
	CreateUser(ctx context.Context, name string, phone string, email string, password string) (string, error)
	ViewUser(ctx context.Context, id string) (string, error)
	UpdateUser(ctx context.Context, id string, name string, phone string, email string) (string, error)
	DeleteUser(ctx context.Context, id string) (string, error)
	ListUser(ctx context.Context, limit string, offset string) (string, error)
	AuthenticateUser(ctx context.Context, phone string, password string) (string, error)
}
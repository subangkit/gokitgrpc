package user

import "context"

type User struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	ViewUser(ctx context.Context, id string) (User, error)
	UpdateUser(ctx context.Context, id string, name string, phone string, email string) error
	DeleteUser(ctx context.Context, id string) error
	ListUser(ctx context.Context, limit string, offset string) (string, error)
	AuthenticateUser(ctx context.Context, phone string, password string) (string, error)
}
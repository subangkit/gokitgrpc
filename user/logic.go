package user

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
)

type service struct {
	repostory Repository
	logger    log.Logger
}

func NewService(rep Repository, logger log.Logger) Service {
	return &service{
		repostory: rep,
		logger:    logger,
	}
}

//CreateUser(ctx context.Context, name string, phone string, email string, password string) (string, error)
func (s service) CreateUser(ctx context.Context, name string, phone string, email string, password string) (string, error) {
	logger := log.With(s.logger, "method", "CreateUser")
	
	uuid, _ := uuid.NewV4()
	id := uuid.String()
	user := User{
		ID:       id,
		Name:     name,
		Phone:    phone,
		Email:    email,
		Password: password,
	}

	if err := s.repostory.CreateUser(ctx, user); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("create user", id)

	return "Success", nil
}
//ViewUser(ctx context.Context, id string) (string, error)
func (s service) ViewUser(ctx context.Context, id string) (string, error) {
	logger := log.With(s.logger, "method", "ViewUser")

	user, err := s.repostory.ViewUser(ctx, id)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	user.Password = ""

	userJson, err := json.Marshal(user)
    if err != nil {
        return "", err
    }

	logger.Log("Get user", id)

	return string(userJson), nil
}
func (s service) UpdateUser(ctx context.Context, id string, name string, phone string, email string) (string, error) {
	logger := log.With(s.logger, "method", "UpdateUser")

	if err := s.repostory.UpdateUser(ctx, id, name, phone, email); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("Update user", id)
	return "Success", nil
}
func (s service) DeleteUser(ctx context.Context, id string) (string, error) {
	logger := log.With(s.logger, "method", "DeleteUser")

	if err := s.repostory.DeleteUser(ctx, id); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("Delete user", id)
	return "Success", nil
}

func (s service) ListUser(ctx context.Context, limit string, offset string) (string, error) {
	logger := log.With(s.logger, "method", "ListUser")
	logger.Log("list user")

	users, err := s.repostory.ListUser(ctx, limit, offset)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	return users, nil
}
func (s service) AuthenticateUser(ctx context.Context, phone string, password string) (string, error) {
	logger := log.With(s.logger, "method", "AuthenticateUser")
	logger.Log("auth user")

	access_token, err := s.repostory.AuthenticateUser(ctx, phone, password)

	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	return access_token, nil
}

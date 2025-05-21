package service

import (
	"context"

	userpb "github.com/vadim8q258475/store-user-microservice/gen/v1"
)

type UserService interface {
	Get(ctx context.Context, userId uint32) (*userpb.User, error)
}

type userService struct {
	client userpb.UserServiceClient
}

func NewUserService(client userpb.UserServiceClient) UserService {
	return &userService{
		client: client,
	}
}

func (s *userService) Get(ctx context.Context, userId uint32) (*userpb.User, error) {
	request := &userpb.GetByID_Request{Id: userId}
	response, err := s.client.GetByID(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.User, nil
}

package grpc_server

import (
	"context"
	"elex_storage/identity_service/internal/use_case/cqrs/command_handlers"
	"elex_storage/identity_service/internal/use_case/cqrs/commands"
	"elex_storage/pkg/shared_kernel/grpc_service"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	grpc_service.UnimplementedUserServiceServer
	registerUserHandler *command_handlers.RegisterUserHandler
	loginUserHandler    *command_handlers.LoginUserHandler
}

func NewUserService(registerUserHandler *command_handlers.RegisterUserHandler, loginUserHandler *command_handlers.LoginUserHandler) grpc_service.UserServiceServer {
	return &UserServer{registerUserHandler: registerUserHandler, loginUserHandler: loginUserHandler}
}

func (s *UserServer) Register(ctx context.Context, req *grpc_service.RegisterUserRequest) (*grpc_service.TokenResponse, error) {
	cmd := commands.RegisterUserCommand{}
	copier.Copy(&cmd, &req)
	err := cmd.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	registerResponse, err := s.registerUserHandler.Handle(cmd)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	var result grpc_service.TokenResponse
	copier.Copy(&result, &registerResponse)
	return &result, nil
}

func (s *UserServer) Login(ctx context.Context, req *grpc_service.LoginUserRequest) (*grpc_service.TokenResponse, error) {
	cmd := commands.LoginUserCommand{}
	copier.Copy(&cmd, &req)
	err := cmd.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	loginResponse, err := s.loginUserHandler.Handle(cmd)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	var result grpc_service.TokenResponse
	copier.Copy(&result, &loginResponse)
	return &result, nil
}

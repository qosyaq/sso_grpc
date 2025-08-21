package auth

import (
	"context"
	"errors"
	"sso/internal/services/auth"
	val "sso/internal/validator"
	"sso/protos/gen/go/proto/sso"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ! Интерфейс сервиса Авторизации.

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth      Auth
	validator *val.Validator
}

func Register(gRPC *grpc.Server, auth Auth, validator *val.Validator) {
	sso.RegisterAuthServer(gRPC, &serverAPI{auth: auth, validator: validator})
}

func (s *serverAPI) Login(ctx context.Context, req *sso.LoginRequest,
) (*sso.LoginResponse, error) {
	verrs, err := s.validator.Auth.ValidateLogin(req)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal validation error")
	}

	if verrs != nil {
		messages := s.validator.HumanizeValidationErrors(verrs)
		return nil, status.Error(codes.InvalidArgument, strings.Join(messages, "; "))
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &sso.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {
	verrs, err := s.validator.Auth.ValidateRegister(req)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal validation error")
	}

	if verrs != nil {
		messages := s.validator.HumanizeValidationErrors(verrs)
		return nil, status.Error(codes.InvalidArgument, strings.Join(messages, "; "))
	}

	uid, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &sso.RegisterResponse{UserId: uid}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {
	verrs, err := s.validator.Auth.ValidateIsAdmin(req)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal validation error")
	}

	if verrs != nil {
		messages := s.validator.HumanizeValidationErrors(verrs)
		return nil, status.Error(codes.InvalidArgument, strings.Join(messages, "; "))
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to check admin status")
	}

	return &sso.IsAdminResponse{IsAdmin: isAdmin}, nil
}

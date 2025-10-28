package grpc

import (
	"context"

	pb "github.com/ploezy/ecommerce-platform/proto/user"
	"github.com/ploezy/ecommerce-platform/user-service/internal/service"
	"github.com/ploezy/ecommerce-platform/user-service/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCServer struct {
	pb.UnimplementedUserServiceServer
	service   service.UserService
	jwtSecret string
}

func NewUserGRPCServer(service service.UserService, jwtSecret string) *UserGRPCServer {
	return &UserGRPCServer{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

func (s *UserGRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := s.service.Register(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.RegisterResponse{
		Id:        uint32(user.ID),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Message:   "User registered successfully",
	}, nil
}

func (s *UserGRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, user, err := s.service.Login(req.Email, req.Password, s.jwtSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
	}

	return &pb.LoginResponse{
		Token:   token,
		UserId:  uint32(user.ID),
		Email:   user.Email,
		Role:    user.Role,
		Message: "Login successful",
	}, nil
}

func (s *UserGRPCServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error) {
	user, err := s.service.GetByID(uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.UserResponse{
		Id:        uint32(user.ID),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}, nil
}

func (s *UserGRPCServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error) {
	user, err := s.service.GetByEmail(req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.UserResponse{
		Id:        uint32(user.ID),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}, nil
}

func (s *UserGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := auth.ValidateToken(req.Token, s.jwtSecret)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: uint32(claims.UserID),
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}
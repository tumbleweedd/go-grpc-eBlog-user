package service

import (
	"context"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/client"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/internal/repository"
	pb2 "github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pkg/pb"
)

type User interface {
	GetUserList(ctx context.Context, req *pb2.GetUserListRequest) (*pb2.GetUserListResponse, error)
	GetLoggedUserProfile(ctx context.Context, req *pb2.GetLoggedUserProfileRequest) (*pb2.GetLoggedUserProfileResponse, error)
	AddUser(ctx context.Context, req *pb2.AddUserRequest) (*pb2.AddUserResponse, error)
	GiveRole(ctx context.Context, req *pb2.GiveRoleRequest) (*pb2.GiveRoleResponse, error)
	GetUserProfile(ctx context.Context, req *pb2.GetUserProfileRequest) (*pb2.GetUserProfileResponse, error)
	DeleteUser(ctx context.Context, req *pb2.DeleteUserRequest) (*pb2.DeleteUserResponse, error)
	UpdateUser(ctx context.Context, req *pb2.UpdateUserRequest) (*pb2.UpdateUserResponse, error)
	GetUserPosts(ctx context.Context, req *pb2.GetUserPostsRequest) (*pb2.GetUserPostsResponse, error)
	GetUserIdByUsername(ctx context.Context, request *pb2.GetUserIdByUsernameRequest) (*pb2.GetUserIdByUsernameResponse, error)
}

type Service struct {
	User
	pb2.UnsafeUserServiceServer
}

func NewService(r *repository.Repository, postSvc client.PostServiceClient) *Service {
	return &Service{
		User: NewUserService(r.User, postSvc),
	}
}

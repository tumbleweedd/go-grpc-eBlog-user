package service

import (
	"context"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/client"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pb"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pkg/repository"
)

type User interface {
	GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListResponse, error)
	GetLoggedUserProfile(ctx context.Context, req *pb.GetLoggedUserProfileRequest) (*pb.GetLoggedUserProfileResponse, error)
	AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error)
	GiveRole(ctx context.Context, req *pb.GiveRoleRequest) (*pb.GiveRoleResponse, error)
	GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error)
	DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	GetUserPosts(ctx context.Context, req *pb.GetUserPostsRequest) (*pb.GetUserPostsResponse, error)
	GetUserIdByUsername(ctx context.Context, request *pb.GetUserIdByUsernameRequest) (*pb.GetUserIdByUsernameResponse, error)
}

type Service struct {
	User
	pb.UnsafeUserServiceServer
}

func NewService(r *repository.Repository, postSvc client.PostServiceClient) *Service {
	return &Service{
		User: NewUserService(r.User, postSvc),
	}
}

package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/client"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pb"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pkg/models"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pkg/repository"
	"net/http"
	"strings"
)

const (
	salt       = "qelwnjgo23ijqpk1"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
)

type UserService struct {
	userRepository repository.User
	postService    client.PostServiceClient
}

func NewUserService(userRepo repository.User, postSvc client.PostServiceClient) *UserService {
	return &UserService{
		userRepository: userRepo,
		postService:    postSvc,
	}
}

func (userService *UserService) GetUserPosts(ctx context.Context, req *pb.GetUserPostsRequest) (*pb.GetUserPostsResponse, error) {

	id, err := userService.userRepository.GetUserIdByUsername(req.GetUsername())
	if err != nil {
		return &pb.GetUserPostsResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	userPosts, err := userService.postService.GetAllPostsByUserId(id)
	if err != nil {
		return &pb.GetUserPostsResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	data := make([]*pb.Post, 0, len(userPosts.GetData()))
	for _, postData := range userPosts.GetData() {
		data = append(data, &pb.Post{
			Body:     postData.Body,
			Head:     postData.Head,
			Category: postData.Category,
			Tags:     postData.GetTags(),
		})
	}

	return &pb.GetUserPostsResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil
}

func (userService *UserService) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	users, err := userService.userRepository.GetUserList()
	if err != nil {
		return &pb.GetUserListResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	data := make([]*pb.UserData, 0, len(users))

	for _, user := range users {
		data = append(data, &pb.UserData{
			Name:     user.Name,
			Lastname: user.Lastname,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return &pb.GetUserListResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil

}

func (userService *UserService) GetLoggedUserProfile(ctx context.Context, req *pb.GetLoggedUserProfileRequest) (*pb.GetLoggedUserProfileResponse, error) {
	user, err := userService.userRepository.GetLoggedUserProfile(int(req.GetUserId()))
	if err != nil {
		return &pb.GetLoggedUserProfileResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	data := &pb.UserData{
		Name:     user.Name,
		Lastname: user.Lastname,
		Username: user.Username,
		Email:    user.Email,
	}

	return &pb.GetLoggedUserProfileResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil

}

func (userService *UserService) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	userId := req.GetCurrentUserId()

	currentUser, err := userService.userRepository.GetLoggedUserProfile(int(userId))
	if err != nil {
		return &pb.AddUserResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	if currentUser.Role != models.ADMIN_ROLE {
		return &pb.AddUserResponse{
			Status: http.StatusForbidden,
			Error:  "access denied",
		}, nil
	}

	addedUser := &models.UserDTOForAdmin{
		Name:               req.Data.GetName(),
		Lastname:           req.Data.GetLastname(),
		Username:           req.Data.GetUsername(),
		Password:           generateHashPassword(req.Data.GetPassword()),
		Email:              req.Data.GetEmail(),
		Role:               models.Role(req.Data.GetRole()),
		IsAccountNonLocked: req.Data.GetIsAccountNonLocked(),
	}

	id, err := userService.userRepository.AddUser(*addedUser)
	if err != nil {
		return &pb.AddUserResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	return &pb.AddUserResponse{
		UserId: int64(id),
		Status: http.StatusOK,
	}, nil
}

func (userService *UserService) GiveRole(ctx context.Context, req *pb.GiveRoleRequest) (*pb.GiveRoleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (userService *UserService) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	user, err := userService.userRepository.GetUserProfile(req.GetUsername())
	if err != nil {
		return &pb.GetUserProfileResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	data := &pb.UserData{
		Name:     user.Name,
		Lastname: user.Lastname,
		Username: user.Username,
		Email:    user.Email,
	}

	return &pb.GetUserProfileResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil
}

func (userService *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	username := req.GetUsername()
	currentUserId := req.GetUserId()

	currentUser, err := userService.userRepository.GetLoggedUserProfile(int(currentUserId))
	if err != nil {
		return &pb.DeleteUserResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	if strings.EqualFold(currentUser.Username, username) || currentUser.Role == models.ADMIN_ROLE {
		err = userService.userRepository.DeleteUser(username)
		if err != nil {
			return &pb.DeleteUserResponse{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}, nil
		}

		return &pb.DeleteUserResponse{
			Status: http.StatusOK,
		}, nil
	}

	return &pb.DeleteUserResponse{
		Status: http.StatusForbidden,
	}, nil
}

func (userService *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	usernameFromRequestParam := req.GetUsername()
	userId := req.GetUserId()

	currentUser, err := userService.userRepository.GetLoggedUserProfile(int(userId))
	if err != nil {
		return &pb.UpdateUserResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}, nil
	}

	if strings.EqualFold(currentUser.Username, usernameFromRequestParam) || currentUser.Role == models.ADMIN_ROLE {
		username := req.Data.GetUsername()
		email := req.Data.GetEmail()
		var password string
		if req.Data.GetPassword() == "" {
			password = currentUser.Password
		} else {
			password = generateHashPassword(req.Data.GetPassword())
		}

		userBeingUpdated := models.UserUpdate{
			Username: &username,
			Email:    &email,
			Password: &password,
		}

		err = userService.userRepository.UpdateUser(usernameFromRequestParam, userBeingUpdated)
		if err != nil {
			return &pb.UpdateUserResponse{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}, nil
		}

		return &pb.UpdateUserResponse{
			Status: http.StatusOK,
		}, nil
	}

	return &pb.UpdateUserResponse{
		Status: http.StatusForbidden,
	}, nil

}

func generateHashPassword(password string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(salt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}

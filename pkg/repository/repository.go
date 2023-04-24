package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pkg/models"
)

type User interface {
	GetUserList() ([]models.UserDTO, error)
	GetLoggedUserProfile(userId int) (models.User, error)
	AddUser(user models.UserDTOForAdmin) (int, error)
	GiveRole(username string) error
	GetUserProfile(username string) (*models.UserDTO, error)
	DeleteUser(username string) error
	UpdateUser(username string, user models.UserUpdate) error
	GetUserIdByUsername(username string) (models.User, error)
}

type Repository struct {
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}

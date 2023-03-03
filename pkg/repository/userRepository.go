package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/tumbleweedd/grpc-eBlog/grpc-eBlog-user/pkg/models"
	"strings"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (userRepo *UserRepository) GetUserList() ([]models.UserDTO, error) {
	var users []models.UserDTO

	query := fmt.Sprintf(`select u.name, u.lastname, u.username, u.email from %s u`, usersTable)

	err := userRepo.db.Select(&users, query)

	return users, err
}

func (userRepo *UserRepository) GetLoggedUserProfile(userId int) (models.User, error) {
	var user models.User

	query := fmt.Sprintf(`select * from %s u where u.user_id = $1`, usersTable)

	err := userRepo.db.Get(&user, query, userId)

	return user, err
}

func (userRepo *UserRepository) AddUser(user models.UserDTOForAdmin) (int, error) {
	var userId int

	query := fmt.Sprintf(`insert into %s (email, is_account_non_locked, lastname, name, password, role, username) 
								values ($1, $2, $3, $4, $5, $6, $7) returning user_id`, usersTable)
	row := userRepo.db.QueryRow(query, user.Email, user.IsAccountNonLocked, user.Lastname, user.Name, user.Password, user.Role, user.Username)

	if err := row.Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}

func (userRepo *UserRepository) GetUserProfile(username string) (*models.UserDTO, error) {
	var user models.UserDTO

	query := fmt.Sprintf(`select u.name, u.lastname, u.username, u.email from %s u where u.username = $1`, usersTable)

	err := userRepo.db.Get(&user, query, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (userRepo *UserRepository) UpdateUser(username string, user models.UserUpdate) error {
	setElements := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	farmingPartsOfTheRequest(user, &setElements, &args, "username", &argId)
	farmingPartsOfTheRequest(user, &setElements, &args, "email", &argId)
	farmingPartsOfTheRequest(user, &setElements, &args, "password", &argId)

	setQuery := strings.Join(setElements, ", ")

	query := fmt.Sprintf(`update %s u set %s where u.username = $%d`, usersTable, setQuery, argId)

	args = append(args, username)

	_, err := userRepo.db.Exec(query, args...)
	return err
}

func (userRepo *UserRepository) DeleteUser(username string) error {
	query := fmt.Sprintf(`delete from %s u where u.username = $1`, usersTable)

	_, err := userRepo.db.Exec(query, username)

	return err
}

func farmingPartsOfTheRequest(user models.UserUpdate, setElements *[]string, args *[]interface{}, element string, argId *int) {
	var requestedElement *string
	if element == "username" {
		requestedElement = user.Username
	} else if element == "email" {
		requestedElement = user.Email
	} else if element == "password" {
		requestedElement = user.Password
	}

	if requestedElement != nil {
		format := element + "=$%d"
		*setElements = append(*setElements, fmt.Sprintf(format, *argId))
		*args = append(*args, requestedElement)
		*argId++
	}
}

func (userRepo *UserRepository) GiveRole(username string) error {
	query := fmt.Sprintf(`update %s u set role = $1 where u.username = $2`, usersTable)

	_, err := userRepo.db.Exec(query, models.ADMIN_ROLE, username)

	return err
}

func (userRepo *UserRepository) GetUserIdByUsername(username string) (int, error) {
	var userId int

	query := fmt.Sprintf(`select u.user_id from %s u where u.username=$1`, usersTable)

	err := userRepo.db.Get(&userId, query, username)

	return userId, err
}

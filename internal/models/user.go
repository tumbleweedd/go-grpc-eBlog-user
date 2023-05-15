package models

type Role string

const (
	ADMIN_ROLE Role = "ADMIN"
)

type User struct {
	Id                 int    `json:"-" db:"user_id"`
	Name               string `json:"name" db:"name"`
	Lastname           string `json:"lastname" db:"lastname"`
	Username           string `json:"username" db:"username"`
	Email              string `json:"email" db:"email"`
	Password           string `json:"password" db:"password"`
	Role               Role   `json:"-" db:"role"`
	IsAccountNonLocked bool   `json:"-" db:"is_account_non_locked"`
}

type UserDTOForAdmin struct {
	Id                 string `json:"-" db:"user_id"`
	Name               string `json:"name" db:"name"`
	Lastname           string `json:"lastname" db:"lastname"`
	Username           string `json:"username" db:"username"`
	Password           string `json:"password" db:"password"`
	Email              string `json:"email" db:"email"`
	Role               Role   `json:"role" db:"role"`
	IsAccountNonLocked bool   `json:"is_account_non_locked" db:"is_account_non_locked"`
}

type UserDTO struct {
	Name     string `json:"name" db:"name"`
	Lastname string `json:"lastname" db:"lastname"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}

type UserUpdate struct {
	Username *string `json:"username" db:"username"`
	Email    *string `json:"email" db:"email"`
	Password *string `json:"password" db:"password"`
}

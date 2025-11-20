package model

import "time"

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type User struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterDTO struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

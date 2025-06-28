package model

import "time"

type UpdateUserRequest struct {
	Password    string `json:"password" binding:"omitempty,required,min=8,max=64"`
	DisplayName string `json:"displayName" binding:"omitempty"`
	Email       string `json:"email" binding:"omitempty,email"`
}

type UserCreateRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=8,max=64"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email" binding:"omitempty,email"`
	Role        string `json:"role" binding:"omitempty,oneof=super_admin admin project_manager reviewer labeler"`
}

type UserCreateResponse struct {
	Username    string `json:"username" binding:"required"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email" binding:"omitempty,email"`
	Role        string `json:"role" binding:"omitempty,oneof=super_admin admin project_manager reviewer labeler"`
}

type User struct {
	ID          int64     `db:"id"`
	Username    string    `db:"username"`
	Password    string    `db:"password_hash"`
	DisplayName string    `db:"display_name"`
	Email       string    `db:"email"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Role        string    `db:"role"`
}

package model

type UserCreateRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=8,max=64"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email" binding:"omitempty,email"`
	Role        string `json:"role" binding:"omitempty,oneof=super_admin admin project_manager reviewer labeler"`
}

type UserCreateResponse struct {
	Username    string `json:"username" binding:"required"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email" binding:"omitempty,email"`
	Role        string `json:"role" binding:"omitempty,oneof=super_admin admin project_manager reviewer labeler"`
}

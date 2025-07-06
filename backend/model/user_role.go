package model

type UserRole string

const (
	RoleSuperAdmin     UserRole = "super_admin"
	RoleAdmin          UserRole = "admin"
	RoleProjectManager UserRole = "project_manager"
	RoleReviewer       UserRole = "reviewer"
	RoleLabeler        UserRole = "labeler"
	RoleUnassigned     UserRole = "unassigned"
)

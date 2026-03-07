package identity

import (
	"time"
)

type User struct {
	ID        int64      `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Role      string     `json:"role"`
	Enabled   bool       `json:"enabled"`
	CreatedAt *time.Time `json:"created_at"`
}

type Role struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
}

type Permission struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
}

type GetUserReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetUserRes struct {
	User *User `json:"user"`
}

type SearchUsersReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Page     int64  `json:"page"`
	Limit    int64  `json:"limit"`
}

type SearchUsersRes struct {
	Users []*User `json:"users"`
	Total int64   `json:"total"`
}

type AddUserReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type AddUserRes struct {
	ID int64 `json:"id"`
}

type UpdateUserReq struct {
	ID       int64   `json:"id" binding:"required"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Role     *string `json:"role"`
}

type UpdateUserRes struct {
	ID int64 `json:"id"`
}

type GetRoleReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetRoleRes struct {
	Role *Role `json:"role"`
}

type SearchRolesReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Page        int64  `json:"page"`
	Limit       int64  `json:"limit"`
}

type SearchRolesRes struct {
	Roles []*Role `json:"roles"`
	Total int64   `json:"total"`
}

type AddRoleReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AddRoleRes struct {
	ID int64 `json:"id"`
}

type UpdateRoleReq struct {
	ID          int64   `json:"id" binding:"required"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type UpdateRoleRes struct {
	ID int64 `json:"id"`
}

type GetPermissionReq struct {
	ID int64 `json:"id" binding:"required"`
}

type GetPermissionRes struct {
	Permission *Permission `json:"permission"`
}

type SearchPermissionsReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Page        int64  `json:"page"`
	Limit       int64  `json:"limit"`
}

type SearchPermissionsRes struct {
	Permissions []*Permission `json:"permissions"`
	Total       int64         `json:"total"`
}

type AddPermissionReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type AddPermissionRes struct {
	ID int64 `json:"id"`
}

type UpdatePermissionReq struct {
	ID          int64   `json:"id" binding:"required"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type UpdatePermissionRes struct {
	ID int64 `json:"id"`
}

// ===================== Authentication API =====================

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRes struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRes struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type MeRes struct {
	User *User `json:"user"`
}


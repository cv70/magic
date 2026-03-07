package dbdao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type IdentityUserModel struct {
	ID        int64     `gorm:"primarykey;column:id"`
	Username  string    `gorm:"column:username;type:varchar(100);not null;uniqueIndex"`
	Email     string    `gorm:"column:email;type:varchar(255);not null"`
	Password  string    `gorm:"column:password;type:varchar(255)"`
	Role      string    `gorm:"column:role;type:varchar(50)"`
	Enabled   bool      `gorm:"column:enabled;default:true"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime:milli"`
}

func (IdentityUserModel) TableName() string {
	return "identity_users"
}

type IdentityUser struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	Role      string
	Enabled   bool
	CreatedAt time.Time
}

func (d *DB) GetIdentityUser(ctx context.Context, id int64) (*IdentityUser, error) {
	var model IdentityUserModel
	result := d.DB().WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &IdentityUser{
		ID:        model.ID,
		Username:  model.Username,
		Email:     model.Email,
		Password:  model.Password,
		Role:      model.Role,
		Enabled:   model.Enabled,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (d *DB) SearchIdentityUsers(ctx context.Context, username, email, role string, page, limit int64) ([]*IdentityUser, int64, error) {
	var models []IdentityUserModel
	q := d.DB().WithContext(ctx)
	if username != "" {
		q = q.Where("username LIKE ?", "%"+username+"%")
	}
	if email != "" {
		q = q.Where("email LIKE ?", "%"+email+"%")
	}
	if role != "" {
		q = q.Where("role = ?", role)
	}
	var total int64
	q.Model(&IdentityUserModel{}).Count(&total)
	result := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	users := make([]*IdentityUser, len(models))
	for i, m := range models {
		users[i] = &IdentityUser{
			ID:        m.ID,
			Username:  m.Username,
			Email:     m.Email,
			Password:  m.Password,
			Role:      m.Role,
			Enabled:   m.Enabled,
			CreatedAt: m.CreatedAt,
		}
	}
	return users, total, nil
}

func (d *DB) CreateIdentityUser(ctx context.Context, user *IdentityUser) (*IdentityUser, error) {
	model := &IdentityUserModel{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
		Enabled:  user.Enabled,
	}
	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}
	return &IdentityUser{
		ID:        model.ID,
		Username:  model.Username,
		Email:     model.Email,
		Password:  model.Password,
		Role:      model.Role,
		Enabled:   model.Enabled,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (d *DB) UpdateIdentityUser(ctx context.Context, user *IdentityUser) (*IdentityUser, error) {
	updates := map[string]interface{}{}
	if user.Username != "" {
		updates["username"] = user.Username
	}
	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.Password != "" {
		updates["password"] = user.Password
	}
	if user.Role != "" {
		updates["role"] = user.Role
	}
	updates["enabled"] = user.Enabled
	result := d.DB().WithContext(ctx).Model(&IdentityUserModel{}).Where("id = ?", user.ID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return d.GetIdentityUser(ctx, user.ID)
}

type RoleModel struct {
	ID          int64     `gorm:"primarykey;column:id"`
	Name        string    `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	Description string    `gorm:"column:description;type:varchar(500)"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime:milli"`
}

func (RoleModel) TableName() string {
	return "identity_roles"
}

type Role struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
}

func (d *DB) GetRole(ctx context.Context, id int64) (*Role, error) {
	var model RoleModel
	result := d.DB().WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("role not found")
		}
		return nil, result.Error
	}
	return &Role{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
	}, nil
}

func (d *DB) SearchRoles(ctx context.Context, name, description string, page, limit int64) ([]*Role, int64, error) {
	var models []RoleModel
	q := d.DB().WithContext(ctx)
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	var total int64
	q.Model(&RoleModel{}).Count(&total)
	result := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	roles := make([]*Role, len(models))
	for i, m := range models {
		roles[i] = &Role{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
			CreatedAt:   m.CreatedAt,
		}
	}
	return roles, total, nil
}

func (d *DB) CreateRole(ctx context.Context, role *Role) (*Role, error) {
	model := &RoleModel{
		Name:        role.Name,
		Description: role.Description,
	}
	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Role{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
	}, nil
}

func (d *DB) UpdateRole(ctx context.Context, role *Role) (*Role, error) {
	updates := map[string]interface{}{}
	if role.Name != "" {
		updates["name"] = role.Name
	}
	if role.Description != "" {
		updates["description"] = role.Description
	}
	result := d.DB().WithContext(ctx).Model(&RoleModel{}).Where("id = ?", role.ID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("role not found")
	}
	return d.GetRole(ctx, role.ID)
}

type PermissionModel struct {
	ID          int64     `gorm:"primarykey;column:id"`
	Name        string    `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	Description string    `gorm:"column:description;type:varchar(500)"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime:milli"`
}

func (PermissionModel) TableName() string {
	return "identity_permissions"
}

type Permission struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
}

func (d *DB) GetPermission(ctx context.Context, id int64) (*Permission, error) {
	var model PermissionModel
	result := d.DB().WithContext(ctx).Where("id = ?", id).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("permission not found")
		}
		return nil, result.Error
	}
	return &Permission{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
	}, nil
}

func (d *DB) SearchPermissions(ctx context.Context, name, description string, page, limit int64) ([]*Permission, int64, error) {
	var models []PermissionModel
	q := d.DB().WithContext(ctx)
	if name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	var total int64
	q.Model(&PermissionModel{}).Count(&total)
	result := q.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	perms := make([]*Permission, len(models))
	for i, m := range models {
		perms[i] = &Permission{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
			CreatedAt:   m.CreatedAt,
		}
	}
	return perms, total, nil
}

func (d *DB) CreatePermission(ctx context.Context, perm *Permission) (*Permission, error) {
	model := &PermissionModel{
		Name:        perm.Name,
		Description: perm.Description,
	}
	result := d.DB().WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Permission{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
	}, nil
}

func (d *DB) UpdatePermission(ctx context.Context, perm *Permission) (*Permission, error) {
	updates := map[string]interface{}{}
	if perm.Name != "" {
		updates["name"] = perm.Name
	}
	if perm.Description != "" {
		updates["description"] = perm.Description
	}
	result := d.DB().WithContext(ctx).Model(&PermissionModel{}).Where("id = ?", perm.ID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("permission not found")
	}
	return d.GetPermission(ctx, perm.ID)
}

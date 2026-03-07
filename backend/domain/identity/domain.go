package identity

import (
	"backend/datasource/dbdao"
	"backend/utils"
	"context"
	"errors"
)

type IdentityDomain struct {
	DB *dbdao.DB
}

func NewIdentityDomain(db *dbdao.DB) *IdentityDomain {
	return &IdentityDomain{DB: db}
}

func (d *IdentityDomain) GetUser(ctx context.Context, id int64) (*User, error) {
	user, err := d.DB.GetIdentityUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		Enabled:   user.Enabled,
		CreatedAt: &user.CreatedAt,
	}, nil
}

func (d *IdentityDomain) SearchUsers(ctx context.Context, username, email, role string, page, limit int64) ([]*User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	users, total, err := d.DB.SearchIdentityUsers(ctx, username, email, role, page, limit)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*User, len(users))
	for i, u := range users {
		result[i] = &User{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Password:  u.Password,
			Role:      u.Role,
			Enabled:   u.Enabled,
			CreatedAt: &u.CreatedAt,
		}
	}
	return result, total, nil
}

func (d *IdentityDomain) AddUser(ctx context.Context, username, email, password, role string) (*User, error) {
	if username == "" || email == "" {
		return nil, errors.New("missing required fields")
	}
	user := &dbdao.IdentityUser{
		Username: username,
		Email:    email,
		Password: password,
		Role:     role,
		Enabled:  true,
	}
	created, err := d.DB.CreateIdentityUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        created.ID,
		Username:  created.Username,
		Email:     created.Email,
		Password:  created.Password,
		Role:      created.Role,
		Enabled:   created.Enabled,
		CreatedAt: &created.CreatedAt,
	}, nil
}

func (d *IdentityDomain) UpdateUser(ctx context.Context, id int64, username, email, password, role *string) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	user := &dbdao.IdentityUser{ID: id}
	if username != nil {
		user.Username = *username
	}
	if email != nil {
		user.Email = *email
	}
	if password != nil {
		user.Password = *password
	}
	if role != nil {
		user.Role = *role
	}
	updated, err := d.DB.UpdateIdentityUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        updated.ID,
		Username:  updated.Username,
		Email:     updated.Email,
		Password:  updated.Password,
		Role:      updated.Role,
		Enabled:   updated.Enabled,
		CreatedAt: &updated.CreatedAt,
	}, nil
}

func (d *IdentityDomain) GetRole(ctx context.Context, id int64) (*Role, error) {
	role, err := d.DB.GetRole(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Role{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   &role.CreatedAt,
	}, nil
}

func (d *IdentityDomain) SearchRoles(ctx context.Context, name, description string, page, limit int64) ([]*Role, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	roles, total, err := d.DB.SearchRoles(ctx, name, description, page, limit)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Role, len(roles))
	for i, r := range roles {
		result[i] = &Role{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			CreatedAt:   &r.CreatedAt,
		}
	}
	return result, total, nil
}

func (d *IdentityDomain) AddRole(ctx context.Context, name, description string) (*Role, error) {
	if name == "" {
		return nil, errors.New("missing required fields")
	}
	role := &dbdao.Role{
		Name:        name,
		Description: description,
	}
	created, err := d.DB.CreateRole(ctx, role)
	if err != nil {
		return nil, err
	}
	return &Role{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		CreatedAt:   &created.CreatedAt,
	}, nil
}

func (d *IdentityDomain) UpdateRole(ctx context.Context, id int64, name, description *string) (*Role, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	role := &dbdao.Role{ID: id}
	if name != nil {
		role.Name = *name
	}
	if description != nil {
		role.Description = *description
	}
	updated, err := d.DB.UpdateRole(ctx, role)
	if err != nil {
		return nil, err
	}
	return &Role{
		ID:          updated.ID,
		Name:        updated.Name,
		Description: updated.Description,
		CreatedAt:   &updated.CreatedAt,
	}, nil
}

func (d *IdentityDomain) GetPermission(ctx context.Context, id int64) (*Permission, error) {
	perm, err := d.DB.GetPermission(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Permission{
		ID:          perm.ID,
		Name:        perm.Name,
		Description: perm.Description,
		CreatedAt:   &perm.CreatedAt,
	}, nil
}

func (d *IdentityDomain) SearchPermissions(ctx context.Context, name, description string, page, limit int64) ([]*Permission, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	perms, total, err := d.DB.SearchPermissions(ctx, name, description, page, limit)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Permission, len(perms))
	for i, p := range perms {
		result[i] = &Permission{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			CreatedAt:   &p.CreatedAt,
		}
	}
	return result, total, nil
}

func (d *IdentityDomain) AddPermission(ctx context.Context, name, description string) (*Permission, error) {
	if name == "" {
		return nil, errors.New("missing required fields")
	}
	perm := &dbdao.Permission{
		Name:        name,
		Description: description,
	}
	created, err := d.DB.CreatePermission(ctx, perm)
	if err != nil {
		return nil, err
	}
	return &Permission{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		CreatedAt:   &created.CreatedAt,
	}, nil
}

func (d *IdentityDomain) UpdatePermission(ctx context.Context, id int64, name, description *string) (*Permission, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	perm := &dbdao.Permission{ID: id}
	if name != nil {
		perm.Name = *name
	}
	if description != nil {
		perm.Description = *description
	}
	updated, err := d.DB.UpdatePermission(ctx, perm)
	if err != nil {
		return nil, err
	}
	return &Permission{
		ID:          updated.ID,
		Name:        updated.Name,
		Description: updated.Description,
		CreatedAt:   &updated.CreatedAt,
	}, nil
}

// ===================== Authentication Methods =====================

// Login 用户登录
func (d *IdentityDomain) Login(ctx context.Context, username, password string) (*User, string, error) {
	// 查询用户
	users, _, err := d.DB.SearchIdentityUsers(ctx, username, "", "", 1, 1)
	if err != nil || len(users) == 0 {
		return nil, "", errors.New("user not found or password incorrect")
	}

	user := users[0]

	// 验证密码
	if !utils.VerifyPassword(user.Password, password) {
		return nil, "", errors.New("user not found or password incorrect")
	}

	// 生成 token
	token := utils.GenerateToken(user.ID, user.Email)

	return &User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  "", // 不返回密码哈希值
		Role:      user.Role,
		Enabled:   user.Enabled,
		CreatedAt: &user.CreatedAt,
	}, token, nil
}

// Register 用户注册
func (d *IdentityDomain) Register(ctx context.Context, username, email, password string) (*User, string, error) {
	// 检查用户名或邮箱是否已存在
	users, _, err := d.DB.SearchIdentityUsers(ctx, username, "", "", 1, 1)
	if err == nil && len(users) > 0 {
		return nil, "", errors.New("username already exists")
	}

	// 创建新用户
	hashedPassword := utils.HashPassword(password)
	id, err := d.DB.AddIdentityUser(ctx, &dbdao.IdentityUser{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Role:     "user",
		Enabled:  true,
	})
	if err != nil {
		return nil, "", err
	}

	// 获取刚创建的用户
	newUser, err := d.DB.GetIdentityUser(ctx, id)
	if err != nil {
		return nil, "", err
	}

	// 生成 token
	token := utils.GenerateToken(newUser.ID, newUser.Email)

	return &User{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Password:  "", // 不返回密码哈希值
		Role:      newUser.Role,
		Enabled:   newUser.Enabled,
		CreatedAt: &newUser.CreatedAt,
	}, token, nil
}

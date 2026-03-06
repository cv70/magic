package identity

import (
	"backend/utils"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func (d *IdentityDomain) ApiGetUser(c *gin.Context) {
	var req GetUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	user, err := d.GetUser(c, req.ID)
	if err != nil {
		utils.RespError(c, 500, "failed to get user")
		return
	}
	utils.RespSuccess(c, GetUserRes{User: user})
}

func (d *IdentityDomain) ApiSearchUsers(c *gin.Context) {
	var req SearchUsersReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}
	users, total, err := d.SearchUsers(c, req.Username, req.Email, req.Role, req.Page, req.Limit)
	if err != nil {
		utils.RespError(c, 500, "failed to search users")
		return
	}
	utils.RespSuccess(c, SearchUsersRes{Users: users, Total: total})
}

func (d *IdentityDomain) ApiAddUser(c *gin.Context) {
	var req AddUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	user, err := d.AddUser(c, req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		klog.Errorf("failed to add user: %v", err)
		utils.RespError(c, 500, "failed to add user")
		return
	}
	utils.RespSuccess(c, AddUserRes{ID: user.ID})
}

func (d *IdentityDomain) ApiUpdateUser(c *gin.Context) {
	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	user, err := d.UpdateUser(c, req.ID, req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		utils.RespError(c, 500, "failed to update user")
		return
	}
	utils.RespSuccess(c, UpdateUserRes{ID: user.ID})
}

func (d *IdentityDomain) ApiGetRole(c *gin.Context) {
	var req GetRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	role, err := d.GetRole(c, req.ID)
	if err != nil {
		utils.RespError(c, 500, "failed to get role")
		return
	}
	utils.RespSuccess(c, GetRoleRes{Role: role})
}

func (d *IdentityDomain) ApiSearchRoles(c *gin.Context) {
	var req SearchRolesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}
	roles, total, err := d.SearchRoles(c, req.Name, req.Description, req.Page, req.Limit)
	if err != nil {
		utils.RespError(c, 500, "failed to search roles")
		return
	}
	utils.RespSuccess(c, SearchRolesRes{Roles: roles, Total: total})
}

func (d *IdentityDomain) ApiAddRole(c *gin.Context) {
	var req AddRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	role, err := d.AddRole(c, req.Name, req.Description)
	if err != nil {
		utils.RespError(c, 500, "failed to add role")
		return
	}
	utils.RespSuccess(c, AddRoleRes{ID: role.ID})
}

func (d *IdentityDomain) ApiUpdateRole(c *gin.Context) {
	var req UpdateRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	role, err := d.UpdateRole(c, req.ID, req.Name, req.Description)
	if err != nil {
		utils.RespError(c, 500, "failed to update role")
		return
	}
	utils.RespSuccess(c, UpdateRoleRes{ID: role.ID})
}

func (d *IdentityDomain) ApiGetPermission(c *gin.Context) {
	var req GetPermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	perm, err := d.GetPermission(c, req.ID)
	if err != nil {
		utils.RespError(c, 500, "failed to get permission")
		return
	}
	utils.RespSuccess(c, GetPermissionRes{Permission: perm})
}

func (d *IdentityDomain) ApiSearchPermissions(c *gin.Context) {
	var req SearchPermissionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}
	perms, total, err := d.SearchPermissions(c, req.Name, req.Description, req.Page, req.Limit)
	if err != nil {
		utils.RespError(c, 500, "failed to search permissions")
		return
	}
	utils.RespSuccess(c, SearchPermissionsRes{Permissions: perms, Total: total})
}

func (d *IdentityDomain) ApiAddPermission(c *gin.Context) {
	var req AddPermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	perm, err := d.AddPermission(c, req.Name, req.Description)
	if err != nil {
		utils.RespError(c, 500, "failed to add permission")
		return
	}
	utils.RespSuccess(c, AddPermissionRes{ID: perm.ID})
}

func (d *IdentityDomain) ApiUpdatePermission(c *gin.Context) {
	var req UpdatePermissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespError(c, 400, "failed to parse body")
		return
	}
	perm, err := d.UpdatePermission(c, req.ID, req.Name, req.Description)
	if err != nil {
		utils.RespError(c, 500, "failed to update permission")
		return
	}
	utils.RespSuccess(c, UpdatePermissionRes{ID: perm.ID})
}

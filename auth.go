package rbac

import (
	"strings"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/models"
	"github.com/ne7ermore/gRBAC/services"
)

func init() {
	s := models.Get()
	services.Build(s)
}

func CreatePermisson(name, des string) (*services.Permission, error) {
	pp := models.Get().GetPermissionPools()
	defer pp.Close()

	return services.CreatePermisson(name, des, pp)
}

func GetPerm(id string) (*services.Permission, error) {
	pp := models.Get().GetPermissionPools()
	defer pp.Close()

	return services.GetPermById(id, pp)
}

func UpdatePerm(id string, params ...string) (*services.Permission, error) {
	if len(params) == 0 || len(params) > 2 {
		GetPerm(id)
	}

	updateParams := map[string]string{}
	if len(params) == 1 {
		updateParams["descrip"] = params[0]
	} else {
		updateParams["descrip"] = params[0]
		updateParams["name"] = params[1]
	}

	pp := models.Get().GetPermissionPools()
	defer pp.Close()

	p, err := services.UpdatePerm(id, updateParams, pp)
	if err != nil {
		return p, err
	}

	if _, err = common.Get().
		ResetPerm(id, params[0]); err != nil {
		return nil, err
	}
	return p, nil
}

func CreateRole(name string) (*services.Role, error) {
	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer rp.Close()
	defer pp.Close()

	return services.CreateRole(name, rp, pp)
}

func GetRole(id string) (*services.Role, error) {
	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer rp.Close()
	defer pp.Close()

	return services.GetRoleById(id, rp, pp)
}

func Assign(rid, pid string) (*services.Role, error) {
	r, err := GetRole(rid)
	if err != nil {
		return nil, err
	}

	_, err = GetPerm(pid)
	if err != nil {
		return nil, err
	}

	if _, ok := r.Permissions[pid]; ok {
		return r, common.ErrPermExist
	}

	var updateParams map[string]string

	// for nil Permissions
	if len(r.Permissions) == 0 {
		updateParams = map[string]string{
			"permissions": pid,
		}
	} else {
		prms := make([]string, 0, len(r.Permissions)+1)
		for p, _ := range r.Permissions {
			prms = append(prms, p)
		}
		prms = append(prms, pid)

		updateParams = map[string]string{
			"permissions": strings.Join(prms, common.MongoRoleSep),
		}
	}

	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer rp.Close()
	defer pp.Close()

	r, err = services.UpdateRole(rid, updateParams, rp, pp)
	if err != nil {
		return nil, err
	}

	if err = services.Assign(rid, pid); err != nil {
		return nil, err
	}

	return r, nil
}

func Revoke(rid, pid string) (*services.Role, error) {
	r, err := GetRole(rid)
	if err != nil {
		return nil, err
	}

	_, err = GetPerm(pid)
	if err != nil {
		return nil, err
	}

	if _, ok := r.Permissions[pid]; !ok {
		return r, common.ErrPermNotExist
	}
	var updateParams map[string]string

	if len(r.Permissions) == 1 {
		updateParams = map[string]string{
			"permissions": "",
		}
	} else {
		prms := make([]string, 0, len(r.Permissions)-1)
		for p, _ := range r.Permissions {
			if p == pid {
				continue
			}
			prms = append(prms, p)
		}

		updateParams = map[string]string{
			"permissions": strings.Join(prms, common.MongoRoleSep),
		}
	}

	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer rp.Close()
	defer pp.Close()

	r, err = services.UpdateRole(rid, updateParams, rp, pp)
	if err != nil {
		return nil, err
	}

	if err = services.Revoke(rid, pid); err != nil {
		return nil, err
	}

	return r, nil
}

func CreateUser(uid string) (*services.User, error) {
	up := models.Get().GetUserPools()
	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer pp.Close()
	defer up.Close()
	defer rp.Close()
	defer pp.Close()

	return services.CreateUser(uid, up, rp, pp)
}

func GetUser(id string) (*services.User, error) {
	up := models.Get().GetUserPools()
	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer up.Close()
	defer rp.Close()

	return services.GetUserById(id, up, rp, pp)
}

func GetUserByUid(uid string) (*services.User, error) {
	up := models.Get().GetUserPools()
	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer pp.Close()
	defer up.Close()
	defer rp.Close()

	return services.GetUserByUid(uid, up, rp, pp)
}

func AddRole(id, rid string) (*services.User, error) {
	u, err := GetUser(id)
	if err != nil {
		return nil, err
	}

	_, err = GetRole(rid)
	if err != nil {
		return nil, err
	}

	if _, ok := u.Roles[rid]; ok {
		return u, common.ErrUserRoleExist
	}

	var updateParams map[string]string

	if len(u.Roles) == 0 {
		updateParams = map[string]string{
			"roles": rid,
		}
	} else {
		rs := make([]string, 0, len(u.Roles)+1)
		for r, _ := range u.Roles {
			rs = append(rs, r)
		}
		rs = append(rs, rid)

		updateParams = map[string]string{
			"roles": strings.Join(rs, common.MongoRoleSep),
		}
	}

	up := models.Get().GetUserPools()
	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer pp.Close()
	defer up.Close()
	defer rp.Close()

	u, err = services.UpdateUser(id, updateParams, up, rp, pp)
	if err != nil {
		return nil, err
	}

	if err = services.AddRole(id, rid); err != nil {
		return nil, err
	}

	return u, nil
}

func DelRole(id, rid string) (*services.User, error) {
	u, err := GetUser(id)
	if err != nil {
		return nil, err
	}

	_, err = GetRole(rid)
	if err != nil {
		return nil, err
	}

	if _, ok := u.Roles[rid]; !ok {
		return u, common.ErrUserNotRoleExist
	}

	var updateParams map[string]string

	if len(u.Roles) == 1 {
		updateParams = map[string]string{
			"roles": "",
		}
	} else {
		rs := make([]string, 0, len(u.Roles)-1)
		for r, _ := range u.Roles {
			if r == rid {
				continue
			}
			rs = append(rs, r)
		}

		updateParams = map[string]string{
			"roles": strings.Join(rs, common.MongoRoleSep),
		}
	}

	up := models.Get().GetUserPools()
	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer pp.Close()
	defer up.Close()
	defer rp.Close()

	u, err = services.UpdateUser(id, updateParams, up, rp, pp)
	if err != nil {
		return nil, err
	}

	if err = services.DelRole(id, rid); err != nil {
		return nil, err
	}

	return u, nil
}

// to check if user owns one permission
func IsPrmitted(id, pid string) (bool, error) {
	return common.Get().Permit(id, pid)
}

// to check if role own one permission or not
func IsRolePermitted(roleid, pid string) (bool, error) {
	return common.Get().RolePermit(roleid, pid)
}

func GetAllPerms(skip, limit int, field string) ([]*services.Permission, error) {
	pp := models.Get().GetPermissionPools()
	defer pp.Close()

	return services.GetPerms(skip, limit, field, pp)
}

func GetAllRoles(skip, limit int, field string) ([]*services.Role, error) {
	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer rp.Close()
	defer pp.Close()

	return services.GetRoles(skip, limit, field, rp, pp)
}

func GetAllUsers(skip, limit int, field string) ([]*services.User, error) {
	up := models.Get().GetUserPools()
	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer pp.Close()
	defer up.Close()
	defer rp.Close()

	return services.GetUsers(skip, limit, field, up, rp, pp)
}

func GetPermByDesc(descrip string) (*services.Permission, error) {
	pp := models.Get().GetPermissionPools()
	defer pp.Close()

	return services.GetPermByDesc(descrip, pp)
}

func GetRoleByName(name string) (*services.Role, error) {
	rp := models.Get().GetRolePools()
	pp := models.Get().GetPermissionPools()
	defer rp.Close()
	defer pp.Close()

	return services.GetRoleByName(name, rp, pp)
}

func GetPermsCount() int {
	pp := models.Get().GetPermissionPools()
	defer pp.Close()

	return services.GetPermissionsCount(pp)
}
func GetRolesCount() int {
	rp := models.Get().GetRolePools()
	defer rp.Close()

	return services.GetRolesCount(rp)
}
func GetUsersCount() int {
	up := models.Get().GetUserPools()
	defer up.Close()

	return services.GetUsersCount(up)
}

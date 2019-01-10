package services

import (
	"strings"
	"time"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/plugin"
)

type perMap map[string]plugin.Permission

type Role struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Permissions perMap    `json:"permissions"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

func NewRoleFromModel(r plugin.Role, pp plugin.PermissionPools) *Role {
	_map := make(perMap)
	for _, p := range strings.Split(r.GetPermissions(), common.MongoRoleSep) {
		// skip empty str while m.Permissions is ""
		if p == "" {
			continue
		}

		if _, found := _map[p]; found {
			continue
		}

		if perm, err := pp.Get(p); err == nil {
			_map[p] = perm
		}
	}

	return &Role{
		Id:          r.Getid(),
		Name:        r.GetName(),
		Permissions: _map,
		CreateTime:  r.GetCreateTime(),
		UpdateTime:  r.GetUpdateTime(),
	}
}

func CreateRole(name string, rp plugin.RolePools, pp plugin.PermissionPools) (*Role, error) {
	id, err := rp.New(name)
	if err != nil {
		return nil, err
	}

	mr, err := GetRoleById(id, rp, pp)
	if err != nil {
		return nil, err
	}
	common.Get().NewRole(common.NewStdRole(id))

	return mr, nil
}

func GetRoleById(id string, rp plugin.RolePools, pp plugin.PermissionPools) (*Role, error) {
	mr, err := rp.Get(id)
	if err != nil {
		return nil, err
	}

	return NewRoleFromModel(mr, pp), nil
}

func UpdateRole(id string, update map[string]string, rp plugin.RolePools, pp plugin.PermissionPools) (*Role, error) {

	if err := rp.Update(id, update); err != nil {
		return nil, err
	}

	return GetRoleById(id, rp, pp)
}

func Assign(rid, pid string) error {
	auth := common.Get()
	p, err := auth.GetPerm(pid)
	if err != nil {
		return err
	}

	r, err := auth.GetRole(rid)
	if err != nil {
		return err
	}

	return auth.Assign(r, p)
}

func Revoke(rid, pid string) error {
	auth := common.Get()
	p, err := auth.GetPerm(pid)
	if err != nil {
		return err
	}

	r, err := auth.GetRole(rid)
	if err != nil {
		return err
	}

	return auth.Revoke(r, p)
}

func GetRoles(skip, limit int, field string, rp plugin.RolePools, pp plugin.PermissionPools) ([]*Role, error) {
	rs, err := rp.Gets(skip, limit, field)
	if err != nil {
		return nil, err
	}

	roles := make([]*Role, 0, limit)
	for _, r := range rs {
		roles = append(roles, NewRoleFromModel(r, pp))
	}

	return roles, nil
}

// get role from db by role name
func GetRoleByName(name string, r plugin.RolePools, pp plugin.PermissionPools) (*Role, error) {
	rr, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}

	return NewRoleFromModel(rr, pp), nil
}

func GetRolesCount(rp plugin.RolePools) int {
	return rp.Counts()
}

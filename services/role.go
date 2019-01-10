package services

import (
	"strings"
	"time"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/plugin"
)

type perMap map[string]interface{}

type Role struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Permissions perMap    `json:"permissions"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

func NewRoleFromModel(r plugin.Role) *Role {
	_map := make(perMap)
	for _, p := range strings.Split(r.GetPermissions(), common.MongoRoleSep) {
		// skip empty str while m.Permissions is ""
		if p == "" {
			continue
		}

		if _, found := _map[p]; found {
			continue
		}
		_map[p] = p
	}

	return &Role{
		Id:          r.Getid(),
		Name:        r.GetName(),
		Permissions: _map,
		CreateTime:  r.GetCreateTime(),
		UpdateTime:  r.GetUpdateTime(),
	}
}

func CreateRole(name string, rp plugin.RolePools) (*Role, error) {
	id, err := rp.New(name)
	if err != nil {
		return nil, err
	}

	mr, err := GetRoleById(id, rp)
	if err != nil {
		return nil, err
	}
	common.Get().NewRole(common.NewStdRole(id))

	return mr, nil
}

func GetRoleById(id string, rp plugin.RolePools) (*Role, error) {
	mr, err := rp.Get(id)
	if err != nil {
		return nil, err
	}

	return NewRoleFromModel(mr), nil
}

func UpdateRole(id string, update map[string]string, rp plugin.RolePools) (*Role, error) {

	if err := rp.Update(id, update); err != nil {
		return nil, err
	}

	return GetRoleById(id, rp)
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

func GetRoles(skip, limit int, field string, rp plugin.RolePools) ([]*Role, error) {
	rs, err := rp.Gets(skip, limit, field)
	if err != nil {
		return nil, err
	}

	roles := make([]*Role, 0, limit)
	for _, r := range rs {
		roles = append(roles, NewRoleFromModel(r))
	}

	return roles, nil
}

// get role from db by role name
func GetRoleByName(name string, r plugin.RolePools) (*Role, error) {
	rr, err := r.GetByName(name)
	if err != nil {
		return nil, err
	}

	return NewRoleFromModel(rr), nil
}

func GetRolesCount(rp plugin.RolePools) int {
	return rp.Counts()
}

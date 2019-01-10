package services

import (
	"time"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/plugin"
)

type Permission struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Descrip    string    `json:"descrip"`
	Sep        string    `json:"sep"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

func NewPermissionFromModel(m plugin.Permission) *Permission {
	return &Permission{
		Id:         m.Getid(),
		Descrip:    m.GetDescrip(),
		Name:       m.GetName(),
		Sep:        m.GetSep(),
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
}

func CreatePermisson(name, des string, p plugin.PermissionPools) (*Permission, error) {
	id, err := p.New(name, des)
	if err != nil {
		return nil, err
	}

	pp, err := GetPermById(id, p)
	if err != nil {
		return nil, err
	}

	common.Get().NewPerm(common.NewFirstP(id, des))

	return pp, nil
}

func GetPermById(id string, p plugin.PermissionPools) (*Permission, error) {
	pp, err := p.Get(id)
	if err != nil {
		return nil, err
	}
	return NewPermissionFromModel(pp), nil
}

func GetPermByDesc(descrip string, p plugin.PermissionPools) (*Permission, error) {
	pp, err := p.GetByDesc(descrip)
	if err != nil {
		return nil, err
	}
	return NewPermissionFromModel(pp), nil
}

func UpdatePerm(id string, update map[string]string, p plugin.PermissionPools) (*Permission, error) {

	if err := p.Update(id, update); err != nil {
		return nil, err
	}

	return GetPermById(id, p)
}

func GetPerms(skip, limit int, field string, p plugin.PermissionPools) ([]*Permission, error) {
	ps, err := p.Gets(skip, limit, field)
	if err != nil {
		return nil, err
	}

	perms := make([]*Permission, 0, limit)
	for _, p := range ps {
		perms = append(perms, NewPermissionFromModel(p))
	}

	return perms, nil
}

func GetPermissionsCount(pp plugin.PermissionPools) int {
	return pp.Counts()
}

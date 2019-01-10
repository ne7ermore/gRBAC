package services

import (
	"strings"
	"time"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/plugin"
)

type roleMap map[string]plugin.Role

type User struct {
	Id         string    `json:"id"`
	UserId     string    `json:"user_id"`
	Roles      roleMap   `json:"roles"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

func NewUserFromModel(m plugin.User, rp plugin.RolePools) *User {
	_map := make(roleMap)
	for _, r := range strings.Split(m.GetRoles(), common.MongoRoleSep) {
		if r == "" {
			continue
		}

		if _, found := _map[r]; found {
			continue
		}

		if role, err := rp.Get(r); err != nil {
			continue
		} else {
			_map[r] = role
		}
	}

	return &User{
		Id:         m.Getid(),
		UserId:     m.GetUserId(),
		Roles:      _map,
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
}

func CreateUser(uid string, up plugin.UserPools, rp plugin.RolePools) (*User, error) {
	id, err := up.New(uid)
	if err != nil {
		return nil, err
	}

	u, err := GetUserById(id, up, rp)
	if err != nil {
		return nil, err
	}
	common.Get().NewUser(id)

	return u, nil
}

func GetUserById(id string, up plugin.UserPools, rp plugin.RolePools) (*User, error) {
	u, err := up.Get(id)
	if err != nil {
		return nil, err
	}

	return NewUserFromModel(u, rp), nil
}

func GetUserByUid(uid string, up plugin.UserPools, rp plugin.RolePools) (*User, error) {
	u, err := up.GetByUid(uid)
	if err != nil {
		return nil, err
	}
	return NewUserFromModel(u, rp), nil
}

func UpdateUser(id string, update map[string]string, up plugin.UserPools, rp plugin.RolePools) (*User, error) {

	if err := up.Update(id, update); err != nil {
		return nil, err
	}

	return GetUserById(id, up, rp)
}

func AddRole(uid, rid string) error {
	return common.Get().AddRole(uid, rid)
}

func DelRole(uid, rid string) error {
	return common.Get().DelRole(uid, rid)
}

func GetUsers(skip, limit int, field string, up plugin.UserPools, rp plugin.RolePools) ([]*User, error) {
	us, err := up.Gets(skip, limit, field)
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0, limit)
	for _, u := range us {
		users = append(users, NewUserFromModel(u, rp))
	}

	return users, nil
}

func GetUsersCount(up plugin.UserPools) int {
	return up.Counts()
}

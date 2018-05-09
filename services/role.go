package services

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/models"
)

type perMap map[string]interface{}

type Role struct {
	Id          bson.ObjectId `json:"id"`
	Name        string        `json:"name"`
	Permissions perMap        `json:"permissions"`
	CreateTime  time.Time     `json:"createTime"`
	UpdateTime  time.Time     `json:"updateTime"`
}

func NewRoleFromModel(m *models.Role) *Role {
	_map := make(perMap)
	for _, p := range strings.Split(m.Permissions, common.MongoRoleSep) {
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
		Id:          m.Id,
		Name:        m.Name,
		Permissions: _map,
		CreateTime:  m.CreateTime,
		UpdateTime:  m.UpdateTime,
	}
}

func CreateRole(name string) (*Role, error) {
	col := models.NewRoleColl()
	defer col.Database.Session.Close()

	id := bson.NewObjectId()
	if err := col.Insert(models.Role{
		Id:         id,
		Name:       name,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}); err != nil {
		return nil, err
	}
	r, err := GetRoleById(id)
	if err != nil {
		return nil, err
	}
	common.Get().NewRole(common.NewStdRole(id.Hex()))

	return r, nil
}

func GetRoleById(id bson.ObjectId) (*Role, error) {
	col := models.NewRoleColl()
	defer col.Database.Session.Close()

	mr := new(models.Role)
	if err := col.FindId(id).One(mr); err != nil {
		return nil, err
	}
	return NewRoleFromModel(mr), nil
}

func UpdateRole(id bson.ObjectId, update bson.M) (*Role, error) {
	col := models.NewRoleColl()
	defer col.Database.Session.Close()

	r := new(models.Role)
	if err := col.FindId(id).One(r); err != nil {
		return nil, err
	}
	update["updateTime"] = time.Now()
	if err := col.UpdateId(id, bson.M{
		"$set": update,
	}); err != nil {
		return nil, err
	}
	return GetRoleById(id)
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

func GetRoles(skip, limit int, field string) ([]*Role, error) {
	col := models.NewRoleColl()
	defer col.Database.Session.Close()

	mr := make([]models.Role, 0, limit)
	if err := col.Find(bson.M{}).Limit(limit).Skip(skip).Sort(field).All(&mr); err != nil {
		return nil, err
	}

	roles := make([]*Role, 0, limit)
	for _, r := range mr {
		roles = append(roles, NewRoleFromModel(&r))
	}

	return roles, nil
}

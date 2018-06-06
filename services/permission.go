package services

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/models"
)

type Permission struct {
	Id         bson.ObjectId `json:"id"`
	Name       string        `json:"name"`
	Descrip    string        `json:"descrip"`
	Sep        string        `json:"sep"`
	CreateTime time.Time     `json:"createTime"`
	UpdateTime time.Time     `json:"updateTime"`
}

func NewPermissionFromModel(m *models.Permission) *Permission {
	return &Permission{
		Id:         m.Id,
		Descrip:    m.Descrip,
		Name:       m.Name,
		Sep:        m.Sep,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
	}
}

func CreatePermisson(name, des string) (*Permission, error) {
	col := models.NewPermissionColl()
	defer col.Database.Session.Close()

	id := bson.NewObjectId()
	if err := col.Insert(models.Permission{
		Id:         id,
		Name:       name,
		Descrip:    des,
		Sep:        common.FirstSep,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}); err != nil {
		return nil, err
	}
	p, err := GetPermById(id)
	if err != nil {
		return nil, err
	}
	common.Get().
		NewPerm(common.NewFirstP(id.Hex(), des))

	return p, nil
}

func GetPermById(id bson.ObjectId) (*Permission, error) {
	col := models.NewPermissionColl()
	defer col.Database.Session.Close()

	mp := new(models.Permission)
	if err := col.Find(bson.M{"_id": id}).One(mp); err != nil {
		return nil, err
	}
	return NewPermissionFromModel(mp), nil
}

func GetPermByDesc(descrip string) (*Permission, error) {
	col := models.NewPermissionColl()
	defer col.Database.Session.Close()

	mp := new(models.Permission)
	if err := col.Find(bson.M{"descrip": descrip}).One(mp); err != nil {
		return nil, err
	}
	return NewPermissionFromModel(mp), nil
}

func UpdatePerm(id bson.ObjectId, update bson.M) (*Permission, error) {
	col := models.NewPermissionColl()
	defer col.Database.Session.Close()

	r := new(models.Permission)
	if err := col.FindId(id).One(r); err != nil {
		return nil, err
	}
	update["updateTime"] = time.Now()
	if err := col.UpdateId(id, bson.M{
		"$set": update,
	}); err != nil {
		return nil, err
	}
	return GetPermById(id)
}

func GetPerms(skip, limit int, field string) ([]*Permission, error) {
	col := models.NewPermissionColl()
	defer col.Database.Session.Close()

	mp := make([]models.Permission, 0, limit)
	if err := col.Find(bson.M{}).Limit(limit).Skip(skip).Sort(field).All(&mp); err != nil {
		return nil, err
	}

	perms := make([]*Permission, 0, limit)
	for _, p := range mp {
		perms = append(perms, NewPermissionFromModel(&p))
	}

	return perms, nil
}

func GetPermsCount() int {
	col := models.NewPermissionColl()
	defer col.Database.Session.Close()

	cnt, _ := col.Count()
	return cnt
}

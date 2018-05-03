package services

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"git.investsavior.com/nccredit/auth/common"
	"git.investsavior.com/nccredit/auth/models"
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
		Sep:        m.Sep,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
	}
}

func CreatePermisson(des string) (*Permission, error) {
	col := models.NewPermissionColl()
	defer col.Database.Session.Close()

	id := bson.NewObjectId()
	if err := col.Insert(models.Permission{
		Id:         id,
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

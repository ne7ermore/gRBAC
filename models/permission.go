package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/plugin"
)

type Permission struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `bson:"name"`
	Descrip    string        `bson:"descrip"`
	Sep        string        `bson:"sep"`
	CreateTime time.Time     `bson:"createTime"`
	UpdateTime time.Time     `bson:"updateTime,omitempty"`
}

func (p Permission) Getid() string {
	return p.Id.Hex()
}

func (p Permission) GetName() string {
	return p.Name
}

func (p Permission) GetDescrip() string {
	return p.Descrip
}

func (p Permission) GetSep() string {
	return p.Sep
}

func (p Permission) GetCreateTime() time.Time {
	return p.CreateTime
}

func (p Permission) GetUpdateTime() time.Time {
	return p.UpdateTime
}

type PermissionColl struct {
	*mgo.Collection
}

func (s *Store) GetPermissionPools() plugin.PermissionPools {
	coll := NewMongodbColl(s.auth, s.perm)

	return plugin.PermissionPools(&PermissionColl{coll})
}

func (p *PermissionColl) Gather() ([]plugin.Permission, error) {
	ps := []*Permission{}

	if err := p.Find(nil).All(&ps); err != nil {
		return nil, err
	}

	pps := make([]plugin.Permission, 0, len(ps))
	for _, p := range ps {
		pps = append(pps, plugin.Permission(p))
	}

	return pps, nil
}

func (p *PermissionColl) New(name, des string) (string, error) {
	id := bson.NewObjectId()
	if err := p.Insert(Permission{
		Id:         id,
		Name:       name,
		Descrip:    des,
		Sep:        common.FirstSep,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}); err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (p *PermissionColl) Get(id string) (plugin.Permission, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, common.ErrInvalidMongoId
	}

	mp := new(Permission)
	if err := p.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(mp); err != nil {
		return nil, err
	}

	return plugin.Permission(mp), nil
}

func (p *PermissionColl) GetByDesc(descrip string) (plugin.Permission, error) {
	mp := new(Permission)
	if err := p.Find(bson.M{"descrip": descrip}).One(mp); err != nil {
		return nil, err
	}

	return plugin.Permission(mp), nil
}

func (p *PermissionColl) Update(id string, update map[string]string) error {
	if !bson.IsObjectIdHex(id) {
		return common.ErrInvalidMongoId
	}

	updateParams := bson.M{}
	for k, v := range update {
		updateParams[k] = v
	}
	updateParams["updateTime"] = time.Now()

	if err := p.UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": updateParams,
	}); err != nil {
		return err
	}

	return nil
}

func (p *PermissionColl) Gets(skip, limit int, field string) ([]plugin.Permission, error) {

	ps := make([]*Permission, 0, limit)
	if err := p.Find(bson.M{}).Limit(limit).Skip(skip).Sort(field).All(&ps); err != nil {
		return nil, err
	}

	pps := make([]plugin.Permission, 0, limit)
	for _, p := range ps {
		pps = append(pps, plugin.Permission(p))
	}

	return pps, nil
}

func (p *PermissionColl) Close() {
	p.Database.Session.Close()
}

func (p *PermissionColl) Counts() int {
	c, _ := p.Count()
	return c
}

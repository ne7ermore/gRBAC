package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/plugin"
)

type Role struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Permissions string        `bson:"permissions"`
	CreateTime  time.Time     `bson:"createTime"`
	UpdateTime  time.Time     `bson:"updateTime,omitempty"`
}

func (r Role) Getid() string {
	return r.Id.Hex()
}

func (r Role) GetName() string {
	return r.Name
}

func (r Role) GetPermissions() string {
	return r.Permissions
}

func (r Role) GetCreateTime() time.Time {
	return r.CreateTime
}

func (r Role) GetUpdateTime() time.Time {
	return r.UpdateTime
}

type RoleColl struct {
	*mgo.Collection
}

func (s *Store) GetRolePools() plugin.RolePools {
	coll := NewMongodbColl(s.auth, s.role)

	return plugin.RolePools(&RoleColl{coll})
}

func (r *RoleColl) Gather() ([]plugin.Role, error) {
	rs := []Role{}
	defer r.Database.Session.Close()

	if err := r.Find(nil).All(&rs); err != nil {
		return nil, err
	}

	prs := make([]plugin.Role, 0, len(rs))
	for _, r := range rs {
		prs = append(prs, plugin.Role(r))
	}

	return prs, nil
}

func (r *RoleColl) New(name string) (string, error) {
	id := bson.NewObjectId()
	if err := r.Insert(Role{
		Id:         id,
		Name:       name,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}); err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (r *RoleColl) Get(id string) (plugin.Role, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, common.ErrInvalidMongoId
	}

	mr := new(Role)
	if err := r.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(mr); err != nil {
		return nil, err
	}

	return plugin.Role(mr), nil
}

func (r *RoleColl) Update(id string, update map[string]string) error {
	if !bson.IsObjectIdHex(id) {
		return common.ErrInvalidMongoId
	}

	updateParams := bson.M{}
	for k, v := range update {
		updateParams[k] = v
	}
	updateParams["updateTime"] = time.Now()

	if err := r.UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": updateParams,
	}); err != nil {
		return err
	}

	return nil
}

func (r *RoleColl) Gets(skip, limit int, field string) ([]plugin.Role, error) {
	rs := make([]*Role, 0, limit)
	if err := r.Find(bson.M{}).Limit(limit).Skip(skip).Sort(field).All(&rs); err != nil {
		return nil, err
	}

	prs := make([]plugin.Role, 0, limit)
	for _, r := range rs {
		prs = append(prs, plugin.Role(r))
	}

	return prs, nil
}

func (r *RoleColl) GetByName(name string) (plugin.Role, error) {
	mr := new(Role)
	if err := r.Find(bson.M{"name": name}).One(mr); err != nil {
		return nil, err
	}

	return plugin.Role(mr), nil
}

func (r *RoleColl) Close() {
	r.Database.Session.Close()
}

func (r *RoleColl) Counts() int {
	c, _ := r.Count()
	return c
}

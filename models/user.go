package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/plugin"
)

type User struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	UserId     string        `bson:"user_id"`
	Roles      string        `bson:"roles"`
	CreateTime time.Time     `bson:"createTime"`
	UpdateTime time.Time     `bson:"updateTime,omitempty"`
}

func (u User) Getid() string {
	return u.Id.Hex()
}

func (u User) GetUserId() string {
	return u.UserId
}

func (u User) GetRoles() string {
	return u.Roles
}

func (u User) GetCreateTime() time.Time {
	return u.CreateTime
}

func (u User) GetUpdateTime() time.Time {
	return u.UpdateTime
}

type UserColl struct {
	*mgo.Collection
}

func (s *Store) GetUserPools() plugin.UserPools {
	coll := NewMongodbColl(s.auth, s.user)

	return plugin.UserPools(&UserColl{coll})
}

func (u *UserColl) Gather() ([]plugin.User, error) {
	us := []User{}

	if err := u.Find(nil).All(&us); err != nil {
		return nil, err
	}

	pus := make([]plugin.User, 0, len(us))
	for _, u := range us {
		pus = append(pus, plugin.User(u))
	}

	return pus, nil
}

func (u *UserColl) New(uid string) (string, error) {
	id := bson.NewObjectId()
	if err := u.Insert(User{
		Id:         id,
		UserId:     uid,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}); err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (u *UserColl) Get(id string) (plugin.User, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, common.ErrInvalidMongoId
	}

	mu := new(User)
	if err := u.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(mu); err != nil {
		return nil, err
	}

	return plugin.User(mu), nil
}

func (u *UserColl) GetByUid(uid string) (plugin.User, error) {
	mu := new(User)
	if err := u.Find(bson.M{"user_id": uid}).One(mu); err != nil {
		return nil, err
	}

	return plugin.User(mu), nil
}

func (u *UserColl) Update(id string, update map[string]string) error {
	if !bson.IsObjectIdHex(id) {
		return common.ErrInvalidMongoId
	}

	updateParams := bson.M{}
	for k, v := range update {
		updateParams[k] = v
	}
	updateParams["updateTime"] = time.Now()

	if err := u.UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": updateParams,
	}); err != nil {
		return err
	}

	return nil
}

func (u *UserColl) Gets(skip, limit int, field string) ([]plugin.User, error) {
	us := make([]*User, 0, limit)
	if err := u.Find(bson.M{}).Limit(limit).Skip(skip).Sort(field).All(&us); err != nil {
		return nil, err
	}

	pus := make([]plugin.User, 0, limit)
	for _, u := range us {
		pus = append(pus, plugin.User(u))
	}

	return pus, nil
}

func (u *UserColl) Close() {
	u.Database.Session.Close()
}

func (u *UserColl) Counts() int {
	c, _ := u.Count()
	return c
}

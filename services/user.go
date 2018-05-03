package services

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/models"
)

type roleMap map[string]interface{}

type User struct {
	Id         bson.ObjectId `json:"id"`
	UserId     string        `json:"user_id"`
	Roles      roleMap       `json:"roles"`
	CreateTime time.Time     `json:"createTime"`
	UpdateTime time.Time     `json:"updateTime"`
}

func NewUserFromModel(m *models.User) *User {
	_map := make(roleMap)
	for _, r := range strings.Split(m.Roles, common.MongoRoleSep) {
		if _, found := _map[r]; found {
			continue
		}
		_map[r] = r
	}

	return &User{
		Id:         m.Id,
		UserId:     m.UserId,
		Roles:      _map,
		CreateTime: m.CreateTime,
		UpdateTime: m.UpdateTime,
	}
}

func CreateUser(uid string) (*User, error) {
	col := models.NewUserColl()
	defer col.Database.Session.Close()

	id := bson.NewObjectId()
	if err := col.Insert(models.User{
		Id:         id,
		UserId:     uid,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}); err != nil {
		return nil, err
	}
	r, err := GetUserById(id)
	if err != nil {
		return nil, err
	}
	common.Get().NewUser(id.Hex())

	return r, nil
}

func GetUserById(id bson.ObjectId) (*User, error) {
	col := models.NewUserColl()
	defer col.Database.Session.Close()

	mu := new(models.User)
	if err := col.FindId(id).One(mu); err != nil {
		return nil, err
	}
	return NewUserFromModel(mu), nil
}

func GetUserByUid(uid string) (*User, error) {
	col := models.NewUserColl()
	defer col.Database.Session.Close()

	mu := new(models.User)
	if err := col.Find(bson.M{"user_id": uid}).One(mu); err != nil {
		return nil, err
	}
	return NewUserFromModel(mu), nil
}

func UpdateUser(id bson.ObjectId, update bson.M) (*User, error) {
	col := models.NewUserColl()
	defer col.Database.Session.Close()

	u := new(models.User)
	if err := col.FindId(id).One(u); err != nil {
		return nil, err
	}
	update["updateTime"] = time.Now()
	if err := col.UpdateId(id, bson.M{
		"$set": update,
	}); err != nil {
		return nil, err
	}
	return GetUserById(id)
}

func AddRole(uid, rid string) error {
	return common.Get().AddRole(uid, rid)
}

func DelRole(uid, rid string) error {
	return common.Get().DelRole(uid, rid)
}

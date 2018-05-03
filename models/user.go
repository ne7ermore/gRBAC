package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	UserId     string        `bson:"user_id"`
	Roles      string        `bson:"roles"`
	CreateTime time.Time     `bson:"createTime"`
	UpdateTime time.Time     `bson:"updateTime,omitempty"`
}

type UserColl struct {
	*mgo.Collection
}

func NewUserColl() *UserColl {
	coll := NewMongodbColl("auth", "user")

	return &UserColl{coll}
}

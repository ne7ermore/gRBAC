package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Role struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Permissions string        `bson:"permissions"`
	CreateTime  time.Time     `bson:"createTime"`
	UpdateTime  time.Time     `bson:"updateTime,omitempty"`
}

type RoleColl struct {
	*mgo.Collection
}

func NewRoleColl() *RoleColl {
	coll := NewMongodbColl("auth", "role")

	return &RoleColl{coll}
}

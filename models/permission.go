package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Permission struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `bson:"name"`
	Descrip    string        `bson:"descrip"`
	Sep        string        `bson:"sep"`
	CreateTime time.Time     `bson:"createTime"`
	UpdateTime time.Time     `bson:"updateTime,omitempty"`
}

type PermissionColl struct {
	*mgo.Collection
}

func NewPermissionColl() *PermissionColl {
	coll := NewMongodbColl("auth", "permission")

	return &PermissionColl{coll}
}

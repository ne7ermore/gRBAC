package models

import (
	"fmt"
	"math"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	Session *mgo.Session
}

var db MongoDB

type MongoInfo struct {
	Addrs     string
	Timeout   int64
	PoolLimit int
}

/**
 * 数据库初始化
 */
func NewMongodb(mi MongoInfo) error {
	info := mgo.DialInfo{
		Addrs:     strings.Split(mi.Addrs, ","),
		Timeout:   time.Duration(mi.Timeout * int64(math.Pow10(9))),
		PoolLimit: mi.PoolLimit,
	}
	// connect db
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return err
	}
	db.Session = session
	fmt.Println("Mongo connected, address: " + mi.Addrs)

	// settings
	db.Session.SetMode(mgo.Strong, true)
	db.Session.SetSocketTimeout(time.Duration(5 * time.Second))
	return nil
}

func GetSession() *mgo.Session {
	return db.Session.Clone()
}

func GlobalSession() *mgo.Session {
	return db.Session
}

func Destroy() {
	if db.Session != nil {
		db.Session.Close()
		db.Session = nil
	}
}

func NewMongodbDB(dbName string) *mgo.Database {
	return GetSession().DB(dbName)
}

func NewMongodbColl(dbName, collName string) *mgo.Collection {
	return NewMongodbDB(dbName).C(collName)
}

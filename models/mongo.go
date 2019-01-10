package models

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	addrs     string = "127.0.0.1:27017"
	auth             = "auth"
	perm             = "permission"
	role             = "role"
	user             = "user"
	timeout   int64  = 5
	poolLimit int    = 1000
)

var (
	mStore     *Store = nil
	mStoreOnce sync.Once
)

type MongoDB struct {
	Session *mgo.Session
}

type Store struct {
	addrs, auth, perm, role, user string
	timeout                       int64
	poolLimit                     int
}

func Get() *Store {
	mStoreOnce.Do(func() {
		mStore = &Store{
			addrs:     addrs,
			auth:      auth,
			perm:      perm,
			role:      role,
			user:      user,
			timeout:   timeout,
			poolLimit: poolLimit,
		}
	})
	return mStore
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
func newMongodb(mi MongoInfo) error {
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

func (s *Store) Build() {
	newMongodb(MongoInfo{
		Addrs:     s.addrs,
		Timeout:   s.timeout,
		PoolLimit: s.poolLimit,
	})
}

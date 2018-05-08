package services

import (
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ne7ermore/gRBAC/models"
)

func Test_create(t *testing.T) {
	models.NewMongodb(models.MongoInfo{"127.0.0.1:27017", 5, 1000})
	p, err := CreatePermisson("oooooo", "ceshiyixia"+time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	if err != nil {
		t.Fatal(err)
	}

	updateParams := bson.M{}
	updateParams["sep"] = "update"
	UpdatePerm(p.Id, updateParams)

	perms, err := GetPerms(1, 5, "-updateTime")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range perms {
		println(p.Id)
	}

	perm, err := GetPermByDesc("abc:vvv:ddd")
	if err != nil {
		t.Fatal(err)
	}
	println(perm.Descrip)
}

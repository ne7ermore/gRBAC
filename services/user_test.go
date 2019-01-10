package services

import (
	"testing"

	"github.com/ne7ermore/gRBAC/models"
	"gopkg.in/mgo.v2/bson"
)

func Test_valid(t *testing.T) {
	models.Get().Build()
	a := "asd"
	if bson.IsObjectIdHex(a) {
		t.Fatal()
	}
	b := bson.NewObjectId()
	c := b.Hex()

	if !bson.IsObjectIdHex(c) {
		t.Fatal()
	}

	users, err := GetUsers(0, 5, "-updateTime", models.Get().GetUserPools(), models.Get().GetRolePools())
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range users {
		println(u.UserId)
	}
}

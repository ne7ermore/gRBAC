package services

import (
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ne7ermore/gRBAC/models"
)

func Test_createrole(t *testing.T) {
	models.NewMongodb(models.MongoInfo{"127.0.0.1:27017", 5, 1000})
	r, err := CreateRole("ceshiyixia" + time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	if err != nil {
		t.Fatal(err)
	}

	updateParams := bson.M{}
	updateParams["permissions"] = "5abb05c9155a5790ddf9656d@@##5abb06b9155a5790ddf9656e"

	_, err = UpdateRole(r.Id, updateParams)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := GetRoles(0, 5, "-updateTime")
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range roles {
		println(r.Id)
	}
}

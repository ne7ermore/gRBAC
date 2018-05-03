package services

import (
	"testing"

	"github.com/ne7ermore/gRBAC/models"
)

func Test_build(t *testing.T) {
	models.NewMongodb(models.MongoInfo{"127.0.0.1:27017", 5, 1000})
	err := InitPerm()
	if err != nil {
		t.Fatal(err)
	}
	err = InitRole()
	if err != nil {
		t.Fatal(err)
	}
}

package services

import (
	"testing"
	// "time"

	"gopkg.in/mgo.v2/bson"
	// "github.com/ne7ermore/gRBAC/models"
)

func Test_valid(t *testing.T) {
	a := "asd"
	if bson.IsObjectIdHex(a) {
		t.Fatal()
	}
	b := bson.NewObjectId()
	c := b.Hex()

	if !bson.IsObjectIdHex(c) {
		t.Fatal()
	}
}

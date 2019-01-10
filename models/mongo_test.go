package models

import (
	"testing"
)

func Test_common(t *testing.T) {
	mi := MongoInfo{"127.0.0.1:27017", 5, 1000}
	err := newMongodb(mi)
	if err != nil {
		println(err.Error())
		t.Fail()
	}
}

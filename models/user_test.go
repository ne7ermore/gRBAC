package models

import "testing"
import "time"

func Test_user(t *testing.T) {
	r := NewUserColl()
	err := r.Insert(User{UserId: "a", CreateTime: time.Now()})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

package models

import "testing"
import "time"

func Test_role(t *testing.T) {
	r := NewRoleColl()
	err := r.Insert(Role{Name: "a", CreateTime: time.Now()})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

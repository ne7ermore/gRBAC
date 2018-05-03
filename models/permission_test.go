package models

import "testing"
import "time"

func Test_Per(t *testing.T) {
	p := NewPermissionColl()
	err := p.Insert(Permission{Descrip: "abc:vvv:ddd", Sep: ":", CreateTime: time.Now()})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

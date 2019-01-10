package models

import "testing"

func Test_Per(t *testing.T) {
	p := Get().GetPermissionPools()
	_, err := p.New("name", "abc:vvv:ddd")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

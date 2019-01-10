package models

import "testing"

func Test_role(t *testing.T) {
	r := Get().GetRolePools()
	_, err := r.New("a")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

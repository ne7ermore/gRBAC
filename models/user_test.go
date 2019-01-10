package models

import "testing"

func Test_user(t *testing.T) {
	r := Get().GetUserPools()
	_, err := r.New("1")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

package common

import (
	"testing"
)

func Test_firstPermission(t *testing.T) {
	profile1 := NewFirstP("1", "profile")
	profile2 := NewFirstP("2", "profile")
	admin := NewFirstP("3", "admin")

	if profile1.match(NewFirstP("5", "std-permission")) {
		t.Fatal("Type assertion issue")
	}

	if !profile1.match(profile1) {
		t.Fatalf("%s should have the permission", profile1.getId())
	}
	if !profile1.match(profile2) {
		t.Fatalf("%s should have the permission", profile1.getId())
	}
	if profile1.match(admin) {
		t.Fatalf("%s should not have the permission", profile1.getId())
	}
}

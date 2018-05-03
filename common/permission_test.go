package common

import (
	"testing"
)

// func TestStdPermission(t *testing.T) {
// 	profile1 := NewStdPermission("profile")
// 	profile2 := NewStdPermission("profile")
// 	admin := NewStdPermission("admin")
// 	if !profile1.match(profile2) {
// 		t.Fatalf("%s should have the permission", profile1.ID())
// 	}
// 	if !profile1.match(profile1) {
// 		t.Fatalf("%s should have the permission", profile1.ID())
// 	}
// 	if profile1.match(admin) {
// 		t.Fatalf("%s should not have the permission", profile1.ID())
// 	}
// 	text, err := json.Marshal(profile1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if string(text) == "\"profile\"" {
// 		t.Fatalf("[\"profile\"] expected, but %s got", text)
// 	}
// 	var p StdPermission
// 	if err := json.Unmarshal(text, &p); err != nil {
// 		t.Fatal(err)
// 	}
// 	if p.ID() != "profile" {
// 		t.Fatalf("[profile] expected, but %s got", p.ID())
// 	}
// }

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

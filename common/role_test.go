package common

import (
	"testing"
)

func Test_StdrA(t *testing.T) {
	rA := NewStdRole("role-a")
	if rA.getId() != "role-a" {
		t.Fatalf("[a] expected, but %s got", rA.getId())
	}
	rA.assign(NewFirstP("permission-a", "permission-a"))
	if !rA.permit(NewFirstP("permission-a", "permission-a")) {
		t.Fatal("[permission-a] should permit to rA")
	}
	if len(rA.getPermissions()) != 1 {
		t.Fatal("[a] should have one permission")
	}

	rA.revoke(NewFirstP("permission-a", "permission-a"))
	if rA.permit(NewFirstP("permission-a", "permission-a")) {
		t.Fatal("[permission-a] should not permit to rA")
	}
	if len(rA.getPermissions()) != 0 {
		t.Fatal("[a] should not have any permission")
	}
}

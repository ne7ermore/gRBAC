package common

import (
	"testing"
)

var (
	rA = NewStdRole("role-a")
	pA = NewFirstP("1", "permission-a")
	rB = NewStdRole("role-b")
	pB = NewFirstP("2", "permission-b")
	rC = NewStdRole("role-c")
	pC = NewFirstP("3", "permission-c")

	auth *Auth
)

func assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestauthPrepare(t *testing.T) {
	auth = Get()
	rA.assign(pA)
	rB.assign(pB)
	rC.assign(pC)
}

func TestauthAdd(t *testing.T) {
	assert(t, auth.Add(rA))
	if err := auth.Add(rA); err != ErrRoleExist {
		t.Error("A role can not be readded")
	}
	assert(t, auth.Add(rB))
	assert(t, auth.Add(rC))
}

func TestauthGetRemove(t *testing.T) {
	assert(t, auth.SetParent("role-c", "role-a"))
	assert(t, auth.SetParent("role-a", "role-b"))
	if r, parents, err := auth.Get("role-a"); err != nil {
		t.Fatal(err)
	} else if r.getId() != "role-a" {
		t.Fatalf("[role-a] does not match %s", r.getId())
	} else if len(parents) != 1 {
		t.Fatal("[role-a] should have one parent")
	}
	assert(t, auth.Remove("role-a"))
	if _, ok := auth.roles["role-a"]; ok {
		t.Fatal("Role removing failed")
	}
	if err := auth.Remove("not-exist"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if r, parents, err := auth.Get("role-a"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	} else if r != nil {
		t.Fatal("The instance of role should be a nil")
	} else if parents != nil {
		t.Fatal("The slice of parents should be a nil")
	}
}

func TestauthParents(t *testing.T) {
	assert(t, auth.SetParent("role-c", "role-b"))
	if _, ok := auth.parents["role-c"]["role-b"]; !ok {
		t.Fatal("Parent binding failed")
	}
	assert(t, auth.RemoveParent("role-c", "role-b"))
	if _, ok := auth.parents["role-c"]["role-b"]; ok {
		t.Fatal("Parent unbinding failed")
	}
	if err := auth.RemoveParent("role-a", "role-b"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if err := auth.RemoveParent("role-b", "role-a"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if err := auth.SetParent("role-a", "role-b"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if err := auth.SetParent("role-c", "role-a"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if err := auth.SetParents("role-a", []string{"role-b"}); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if err := auth.SetParents("role-c", []string{"role-a"}); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	assert(t, auth.SetParents("role-c", []string{"role-b"}))
	if _, ok := auth.parents["role-c"]["role-b"]; !ok {
		t.Fatal("Parent binding failed")
	}
	if parents, err := auth.GetParents("role-a"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	} else if len(parents) != 0 {
		t.Fatal("[role-a] should not have any parent")
	}
	if parents, err := auth.GetParents("role-b"); err != nil {
		t.Fatal(err)
	} else if len(parents) != 0 {
		t.Fatal("[role-b] should not have any parent")
	}
	if parents, err := auth.GetParents("role-c"); err != nil {
		t.Fatal(err)
	} else if len(parents) != 1 {
		t.Fatal("[role-c] should have one parent")
	}
}

func TestauthPermission(t *testing.T) {
	if !auth.IsGranted("role-c", pC, nil) {
		t.Fatalf("role-c should have %s", pC)
	}
	if auth.IsGranted("role-c", pC, func(*Auth, string, Permission) bool { return false }) {
		t.Fatal("Assertion don't work")
	}
	if !auth.IsGranted("role-c", pB, nil) {
		t.Fatalf("role-c should have %s which inherits from role-b", pB)
	}

	assert(t, auth.RemoveParent("role-c", "role-b"))
	if auth.IsGranted("role-c", pB, nil) {
		t.Fatalf("role-c should not have %s because of the unbinding with role-b", pB)
	}
}

func Test_A(t *testing.T) {
	a := Get()
	a.NewPerm(pA)
	a.NewRole(rA)
	err := a.Assign(rA, pA)
	if err != nil {
		t.Fatal(err)
	}
	err = a.Assign(rB, pA)
	if err == nil {
		t.Fatal()
	}
	p, err := a.GetPerm("1")
	if err != nil {
		t.Fatal(err)
	}
	println(p.getDes())
}

func BenchmarkauthGranted(b *testing.B) {
	auth = Get()
	rA.assign(pA)
	rB.assign(pB)
	rC.assign(pC)
	auth.Add(rA)
	auth.Add(rB)
	auth.Add(rC)
	for i := 0; i < b.N; i++ {
		auth.IsGranted("role-a", pA, nil)
	}
}

func BenchmarkauthNotGranted(b *testing.B) {
	auth = Get()
	rA.assign(pA)
	rB.assign(pB)
	rC.assign(pC)
	auth.Add(rA)
	auth.Add(rB)
	auth.Add(rC)
	for i := 0; i < b.N; i++ {
		auth.IsGranted("role-a", pB, nil)
	}
}

package common

import (
	"testing"
)

var (
	rA = NewStdRole("role-a")
	pA = NewFirstP("dsa21`212e", "permission-a")
	rB = NewStdRole("role-b")
	pB = NewFirstP("asd2324323432", "permission-b")
	rC = NewStdRole("role-c")
	pC = NewFirstP("12312dsfsd", "permission-c")

	auth *Auth
)

func TestauthPrepare(t *testing.T) {
	auth = Get()
	rA.assign(pA)
	rB.assign(pB)
	rC.assign(pC)
	auth.Assign(rA, pA)
	auth.Assign(rB, pB)
	auth.Assign(rC, pC)
}

func Test_A(t *testing.T) {
	a := Get()
	a.NewPerm(pA)
	a.NewPerm(pB)
	a.NewPerm(pC)
	a.NewRole(rA)
	err := a.Assign(rA, pA)
	if err != nil {
		t.Fatal(err)
	}
	err = a.Assign(rB, pA)
	if err == nil {
		t.Fatal()
	}
	p, err := a.GetPerm("asd2324323432")
	if err != nil {
		t.Fatal(err)
	}
	println(p.getDes())
}

func Test_all(t *testing.T) {
	a := Get()
	// ps := a.GetAllPerms()
	// for _, p := range ps {
	// 	println(p)
	// }
	for k, v := range a.permissions {
		println("=============")
		println(k)
		println(v.getId())
		println(v.getDes())
	}
}

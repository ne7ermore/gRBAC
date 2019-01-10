package services

import (
	"testing"
	"time"

	"github.com/ne7ermore/gRBAC/models"
)

func Test_create(t *testing.T) {
	models.Get().Build()
	p, err := CreatePermisson("oooooo",
		"ceshiyixia"+time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"),
		models.Get().GetPermissionPools())
	if err != nil {
		t.Fatal(err)
	}

	updateParams := map[string]string{"sep": "update"}
	UpdatePerm(p.Id, updateParams, models.Get().GetPermissionPools())

	perms, err := GetPerms(1, 5, "-updateTime", models.Get().GetPermissionPools())
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range perms {
		println(p.Id)
	}

	perm, err := GetPermByDesc("abc:vvv:ddd", models.Get().GetPermissionPools())
	if err != nil {
		t.Fatal(err)
	}
	println(perm.Descrip)
}

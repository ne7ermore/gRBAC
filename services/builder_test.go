package services

import (
	"testing"

	"github.com/ne7ermore/gRBAC/models"
)

func Test_build(t *testing.T) {
	models.Get().Build()
	p := models.Get().GetPermissionPools()
	err := initPerm(p)
	if err != nil {
		t.Fatal(err)
	}

	r := models.Get().GetRolePools()
	err = initRole(r)
	if err != nil {
		t.Fatal(err)
	}
}

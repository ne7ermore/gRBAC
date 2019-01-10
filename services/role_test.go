package services

import (
	"testing"
	"time"

	"github.com/ne7ermore/gRBAC/models"
)

func Test_createrole(t *testing.T) {
	models.Get().Build()
	r, err := CreateRole("ceshiyixia"+time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"), models.Get().GetRolePools(), models.Get().GetPermissionPools())
	if err != nil {
		t.Fatal(err)
	}

	updateParams := map[string]string{"permissions": "5abb05c9155a5790ddf9656d@@##5abb06b9155a5790ddf9656e"}

	_, err = UpdateRole(r.Id, updateParams, models.Get().GetRolePools(), models.Get().GetPermissionPools())
	if err != nil {
		t.Fatal(err)
	}

	roles, err := GetRoles(0, 5, "-updateTime", models.Get().GetRolePools(), models.Get().GetPermissionPools())
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range roles {
		println(r.Id)
	}
}

package services

import (
	"strings"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/plugin"
)

func Build(s plugin.Store) {
	s.Build()

	// init permissions
	p := s.GetPermissionPools()
	if err := initPerm(p); err != nil {
		panic(err)
	}

	// init roles
	r := s.GetRolePools()
	if err := initRole(r); err != nil {
		panic(err)
	}

	// init users
	u := s.GetUserPools()
	if err := initUser(u); err != nil {
		panic(err)
	}
}

func initPerm(pp plugin.PermissionPools) error {
	ps, err := pp.Gather()
	if err != nil {
		return err
	}

	cps := make([]common.Permission, 0, len(ps))
	for _, p := range ps {
		cps = append(cps, common.NewFirstP(p.Getid(), p.GetDescrip()))
	}

	common.Get().NewPerms(cps)
	return nil
}

func initRole(rp plugin.RolePools) error {
	rs, err := rp.Gather()
	if err != nil {
		return err
	}

	auth := common.Get()
	for _, r := range rs {
		sr := common.NewStdRole(r.Getid())
		auth.NewRole(sr)
		if len(r.GetPermissions()) == 0 {
			continue
		}

		permids := strings.Split(r.GetPermissions(), common.MongoRoleSep)
		perms, err := auth.GetPerms(permids)
		if err != nil {
			return err
		}
		for _, p := range perms {
			err = auth.Assign(sr, p)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func initUser(up plugin.UserPools) error {
	us, err := up.Gather()
	if err != nil {
		return err
	}

	auth := common.Get()

	for _, u := range us {
		auth.NewUser(u.Getid())

		if len(u.GetRoles()) == 0 {
			continue
		}

		for _, rid := range strings.Split(u.GetRoles(), common.MongoRoleSep) {
			err = auth.AddRole(u.Getid(), rid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

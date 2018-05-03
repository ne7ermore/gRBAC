package services

import (
	"strings"

	"git.investsavior.com/nccredit/auth/common"
	"git.investsavior.com/nccredit/auth/models"
)

func InitPerm() error {
	mps := []models.Permission{}
	pc := models.NewPermissionColl()
	defer pc.Database.Session.Close()

	if err := pc.Find(nil).All(&mps); err != nil {
		return err
	}

	ps := make([]common.Permission, 0, len(mps))
	for _, mp := range mps {
		ps = append(ps, common.NewFirstP(mp.Id.Hex(), mp.Descrip))
	}
	common.Get().NewPerms(ps)
	return nil
}

func InitRole() error {
	mrs := []models.Role{}
	rc := models.NewRoleColl()
	defer rc.Database.Session.Close()

	if err := rc.Find(nil).All(&mrs); err != nil {
		return err
	}
	auth := common.Get()
	for _, mr := range mrs {
		sr := common.NewStdRole(mr.Id.Hex())
		auth.NewRole(sr)
		if len(mr.Permissions) == 0 {
			continue
		}

		permids := strings.Split(mr.Permissions, common.MongoRoleSep)
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

func InitUser() error {
	mus := []models.User{}
	uc := models.NewUserColl()
	defer uc.Database.Session.Close()

	var (
		uid string
		err error
	)

	if err = uc.Find(nil).All(&mus); err != nil {
		return err
	}
	auth := common.Get()

	for _, mu := range mus {
		uid = mu.Id.Hex()
		auth.NewUser(uid)

		if len(mu.Roles) == 0 {
			continue
		}

		for _, rid := range strings.Split(mu.Roles, common.MongoRoleSep) {
			err = auth.AddRole(uid, rid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

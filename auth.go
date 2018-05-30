package autht

import (
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/ne7ermore/gRBAC/common"
	"github.com/ne7ermore/gRBAC/models"
	"github.com/ne7ermore/gRBAC/services"
)

func init() {
	// init mongodb
	models.NewMongodb(models.MongoInfo{
		Addrs:     common.Addrs,
		Timeout:   common.Timeout,
		PoolLimit: common.PoolLimit,
	})

	// init permissions
	if err := services.InitPerm(); err != nil {
		panic(err)
	}

	// init roles
	if err := services.InitRole(); err != nil {
		panic(err)
	}

	// init users
	if err := services.InitUser(); err != nil {
		panic(err)
	}
}

func CreatePermisson(name, des string) (*services.Permission, error) {
	return services.CreatePermisson(name, des)
}

func GetPerm(id string) (*services.Permission, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, common.ErrInvalidMongoId
	}

	return services.GetPermById(bson.ObjectIdHex(id))
}

func UpdatePerm(id string, params ...string) (*services.Permission, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, common.ErrInvalidMongoId
	}

	if len(params) == 0 || len(params) > 2 {
		GetPerm(id)
	}

	updateParams := bson.M{}
	if len(params) == 1 {
		updateParams["descrip"] = params[0]
	} else {
		updateParams["descrip"] = params[0]
		updateParams["name"] = params[1]
	}

	p, err := services.UpdatePerm(bson.ObjectIdHex(id), updateParams)
	if err != nil {
		return p, err
	}

	if _, err = common.Get().
		ResetPerm(id, params[0]); err != nil {
		return nil, err
	}
	return p, nil
}

func CreateRole(name string) (*services.Role, error) {
	return services.CreateRole(name)
}

func GetRole(id string) (*services.Role, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, common.ErrInvalidMongoId
	}
	return services.GetRoleById(bson.ObjectIdHex(id))
}

func Assign(rid, pid string) (*services.Role, error) {
	if !bson.IsObjectIdHex(rid) {
		return nil, common.ErrInvalidMongoId
	}

	r, err := GetRole(rid)
	if err != nil {
		return nil, err
	}

	_, err = GetPerm(pid)
	if err != nil {
		return nil, err
	}

	if _, ok := r.Permissions[pid]; ok {
		return r, common.ErrPermExist
	}

	var updateParams bson.M

	// for nil Permissions
	if len(r.Permissions) == 0 {
		updateParams = bson.M{
			"permissions": pid,
		}
	} else {
		prms := make([]string, 0, len(r.Permissions)+1)
		for p, _ := range r.Permissions {
			prms = append(prms, p)
		}
		prms = append(prms, pid)

		updateParams = bson.M{
			"permissions": strings.Join(prms, common.MongoRoleSep),
		}
	}

	r, err = services.UpdateRole(bson.ObjectIdHex(rid), updateParams)
	if err != nil {
		return nil, err
	}

	if err = services.Assign(rid, pid); err != nil {
		return nil, err
	}

	return r, nil
}

func Revoke(rid, pid string) (*services.Role, error) {
	if !bson.IsObjectIdHex(rid) {
		return nil, common.ErrInvalidMongoId
	}

	r, err := GetRole(rid)
	if err != nil {
		return nil, err
	}

	_, err = GetPerm(pid)
	if err != nil {
		return nil, err
	}

	if _, ok := r.Permissions[pid]; !ok {
		return r, common.ErrPermNotExist
	}
	var updateParams bson.M

	if len(r.Permissions) == 1 {
		updateParams = bson.M{
			"permissions": "",
		}
	} else {
		prms := make([]string, 0, len(r.Permissions)-1)
		for p, _ := range r.Permissions {
			if p == pid {
				continue
			}
			prms = append(prms, p)
		}

		updateParams = bson.M{
			"permissions": strings.Join(prms, common.MongoRoleSep),
		}
	}

	r, err = services.UpdateRole(bson.ObjectIdHex(rid), updateParams)
	if err != nil {
		return nil, err
	}

	if err = services.Revoke(rid, pid); err != nil {
		return nil, err
	}

	return r, nil
}

func CreateUser(uid string) (*services.User, error) {
	return services.CreateUser(uid)
}

func GetUser(mongoid string) (*services.User, error) {
	if !bson.IsObjectIdHex(mongoid) {
		return nil, common.ErrInvalidMongoId
	}

	return services.GetUserById(bson.ObjectIdHex(mongoid))
}

func GetUserByUid(uid string) (*services.User, error) {
	return services.GetUserByUid(uid)
}

func AddRole(mongoid, rid string) (*services.User, error) {
	if !bson.IsObjectIdHex(mongoid) {
		return nil, common.ErrInvalidMongoId
	}

	u, err := GetUser(mongoid)
	if err != nil {
		return nil, err
	}

	_, err = GetRole(rid)
	if err != nil {
		return nil, err
	}

	if _, ok := u.Roles[rid]; ok {
		return u, common.ErrUserRoleExist
	}

	var updateParams bson.M

	if len(u.Roles) == 0 {
		updateParams = bson.M{
			"roles": rid,
		}
	} else {
		rs := make([]string, 0, len(u.Roles)+1)
		for r, _ := range u.Roles {
			rs = append(rs, r)
		}
		rs = append(rs, rid)

		updateParams = bson.M{
			"roles": strings.Join(rs, common.MongoRoleSep),
		}
	}
	u, err = services.UpdateUser(bson.ObjectIdHex(mongoid), updateParams)
	if err != nil {
		return nil, err
	}

	if err = services.AddRole(mongoid, rid); err != nil {
		return nil, err
	}

	return u, nil
}

func DelRole(mongoid, rid string) (*services.User, error) {
	if !bson.IsObjectIdHex(mongoid) {
		return nil, common.ErrInvalidMongoId
	}

	u, err := GetUser(mongoid)
	if err != nil {
		return nil, err
	}

	_, err = GetRole(rid)
	if err != nil {
		return nil, err
	}

	if _, ok := u.Roles[rid]; !ok {
		return u, common.ErrUserNotRoleExist
	}

	var updateParams bson.M

	if len(u.Roles) == 1 {
		updateParams = bson.M{
			"roles": "",
		}
	} else {
		rs := make([]string, 0, len(u.Roles)-1)
		for r, _ := range u.Roles {
			if r == rid {
				continue
			}
			rs = append(rs, r)
		}

		updateParams = bson.M{
			"roles": strings.Join(rs, common.MongoRoleSep),
		}
	}

	u, err = services.UpdateUser(bson.ObjectIdHex(mongoid), updateParams)
	if err != nil {
		return nil, err
	}

	if err = services.DelRole(mongoid, rid); err != nil {
		return nil, err
	}

	return u, nil
}

// to check if user owns one permission
func IsPrmitted(mongoid, pid string) (bool, error) {
	return common.Get().
		Permit(mongoid, pid)
}

// to check if role own one permission or not
func IsRolePermitted(roleid, pid string) (bool, error) {
	return common.Get().
		RolePermit(roleid, pid)
}

func GetAllPerms(skip, limit int, field string) ([]*services.Permission, error) {
	return services.GetPerms(skip, limit, field)
}

func GetAllRoles(skip, limit int, field string) ([]*services.Role, error) {
	return services.GetRoles(skip, limit, field)
}

func GetAllUsers(skip, limit int, field string) ([]*services.User, error) {
	return services.GetUsers(skip, limit, field)
}

func GetPermByDesc(descrip string) (*services.Permission, error) {
	return services.GetPermByDesc(descrip)
}

func GetRoleByName(name string) (*services.Role, error) {
	return services.GetRoleByName(name)
}

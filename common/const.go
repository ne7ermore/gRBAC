package common

import (
	"errors"
)

var (
	ErrRoleNotExist     = errors.New("Role does not exist")
	ErrRoleTypeNotExist = errors.New("Role type does not exist")
	ErrRoleExist        = errors.New("Role existed")
	ErrPermNotExist     = errors.New("Permission does not exist")
	ErrPermTypeNotExist = errors.New("Permission type does not exist")
	ErrPermExist        = errors.New("Permission existed")
	ErrUserNotExist     = errors.New("User does not exist")
	ErrUserExist        = errors.New("User existed")
	ErrUserRoleExist    = errors.New("User Has role")
	ErrUserNotRoleExist = errors.New("User Has not role")
	ErrInvalidMongoId   = errors.New("Invalid mongo id")
)

var (
	FirstSep     string = ":"
	MongoRoleSep        = "@@##"
)

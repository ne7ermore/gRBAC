package plugin

type Store interface {
	Build()
	GetUserPools() UserPools
	GetPermissionPools() PermissionPools
	GetRolePools() RolePools
}

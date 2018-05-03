package common

type Role interface {
	getId() string
	permit(Permission) bool
}

type Roles map[string]Role
type User map[string]Role

type Permission interface {
	getDes() string
	match(Permission) bool
	getId() string
}

type Permissions map[string]Permission

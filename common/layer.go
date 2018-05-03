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

// type StdPermission struct {
// 	IDStr string
// }

// func NewStdPermission(id string) Permission {
// 	return &StdPermission{id}
// }

// func (p *StdPermission) ID() string {
// 	return p.IDStr
// }

// func (p *StdPermission) Match(a Permission) bool {
// 	return p.IDStr == a.ID()
// }

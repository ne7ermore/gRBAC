package common

import (
	"sync"
)

func NewStdRole(id string) *stdRole {
	return &stdRole{
		id:          id,
		permissions: make(Permissions),
	}
}

type stdRole struct {
	sync.RWMutex
	id          string `json:"id"`
	permissions Permissions
}

func (role *stdRole) getId() string {
	return role.id
}

func (role *stdRole) assign(p Permission) {
	role.Lock()
	role.permissions[p.getId()] = p
	role.Unlock()
}

func (role *stdRole) permit(p Permission) (rslt bool) {
	role.RLock()
	defer role.RUnlock()
	for _, rp := range role.permissions {
		if rp.match(p) {
			rslt = true
			break
		}
	}
	return
}

func (role *stdRole) revoke(p Permission) {
	role.Lock()
	delete(role.permissions, p.getId())
	role.Unlock()
}

func (role *stdRole) getPermissions() []Permission {
	role.RLock()
	defer role.RUnlock()
	result := make([]Permission, 0, len(role.permissions))
	for _, p := range role.permissions {
		result = append(result, p)
	}
	return result
}

func rolesPermit(u User, p Permission) bool {
	if len(u) == 0 {
		return false
	}

	have := make(chan bool)

	for _, r := range u {
		go func(r Role, p Permission) {
			if r.permit(p) {
				have <- true
			} else {
				have <- false
			}
			close(have)
		}(r, p)
	}

	for is := range have {
		if is {
			return true
		}
	}

	return false
}

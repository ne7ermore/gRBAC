package common

import (
	"sync"
)

var (
	authT    *Auth = nil
	authOnce sync.Once
)

type AssertionFunc func(*Auth, string, Permission) bool

type Auth struct {
	sync.RWMutex
	roles       Roles
	parents     map[string]map[string]struct{}
	permissions Permissions
	users       map[string]User
}

func Get() *Auth {
	authOnce.Do(func() {
		authT = &Auth{
			roles:       make(Roles),
			parents:     make(map[string]map[string]struct{}),
			permissions: make(Permissions),
			users:       make(map[string]User),
		}
	})
	return authT
}

func (auth *Auth) NewPerm(p Permission) {
	if fp, ok := p.(*firstPermission); ok {
		auth.Lock()
		auth.permissions[fp.getId()] = fp
		auth.Unlock()
	}
}

func (auth *Auth) NewRole(r Role) {
	if sr, ok := r.(*stdRole); ok {
		auth.Lock()
		auth.roles[sr.getId()] = sr
		auth.Unlock()
	}
}

func (auth *Auth) NewUser(id string) {
	auth.Lock()
	auth.users[id] = make(User)
	auth.Unlock()
}

func (auth *Auth) NewPerms(ps []Permission) {
	auth.Lock()
	for _, p := range ps {
		auth.permissions[p.getId()] = p
	}
	auth.Unlock()
}

func (auth *Auth) Assign(role Role, p Permission) error {
	auth.Lock()
	defer auth.Unlock()

	role, ok := auth.roles[role.getId()]
	if !ok {
		return ErrRoleNotExist
	}
	if sr, ok := role.(*stdRole); ok {
		sr.assign(p)
		return nil
	}
	return ErrRoleTypeNotExist
}

func (auth *Auth) Revoke(role Role, p Permission) error {
	auth.Lock()
	defer auth.Unlock()
	role, ok := auth.roles[role.getId()]
	if !ok {
		return ErrRoleNotExist
	}
	if sr, ok := role.(*stdRole); ok {
		sr.revoke(p)
		return nil
	}
	return ErrRoleTypeNotExist
}

func (auth *Auth) GetAllPerms() Permissions {
	return auth.permissions
}

func (auth *Auth) GetPerms(ids []string) ([]Permission, error) {
	ps := make([]Permission, 0, len(ids))
	for _, id := range ids {
		p, err := auth.GetPerm(id)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (auth *Auth) ResetPerm(id, des string) (Permission, error) {
	auth.RLock()
	defer auth.RUnlock()

	p, err := auth.GetPerm(id)
	if err != nil {
		return p, err
	}

	if fp, ok := p.(*firstPermission); ok {
		fp.setDes(des)
		return fp, nil
	}

	return nil, ErrPermTypeNotExist
}

func (auth *Auth) GetPerm(id string) (Permission, error) {
	auth.RLock()
	defer auth.RUnlock()

	if p, ok := auth.permissions[id]; !ok {
		return nil, ErrPermNotExist
	} else {
		return p, nil
	}
}

func (auth *Auth) GetRole(id string) (Role, error) {
	auth.RLock()
	defer auth.RUnlock()

	if r, ok := auth.roles[id]; !ok {
		return nil, ErrRoleNotExist
	} else {
		return r, nil
	}
}

func (auth *Auth) GetRoles(ids []string) ([]Role, error) {
	rs := make([]Role, 0, len(ids))
	for _, id := range ids {
		r, err := auth.GetRole(id)
		if err != nil {
			return nil, err
		}
		rs = append(rs, r)
	}
	return rs, nil
}

func (auth *Auth) GetAllRoles() Roles {
	return auth.roles
}

// TODO: 解决死锁问题
//
// If a goroutine holds a RWMutex for reading and another goroutine might call Lock,
// no goroutine should expect to be able to acquire a read lock until the initial read lock is released.
// In particular, this prohibits recursive read locking.
// This is to ensure that the lock eventually becomes available;
// a blocked Lock call excludes new readers from acquiring the lock.
//
func (auth *Auth) AddRole(uid, rid string) error {
	// auth.Lock()
	// defer auth.Unlock()

	auth.RLock()
	u, have := auth.users[uid]
	auth.RUnlock()

	if !have {
		return ErrUserNotExist
	}

	r, err := auth.GetRole(rid)
	if err != nil {
		return err
	}

	if _, have = u[rid]; have {
		return ErrUserRoleExist
	}

	// u[rid] = r

	auth.Lock()
	u[rid] = r
	auth.Unlock()

	return nil
}

func (auth *Auth) DelRole(uid, rid string) error {
	auth.Lock()
	defer auth.Unlock()

	u, have := auth.users[uid]
	if !have {
		return ErrUserNotExist
	}

	if _, have = u[rid]; !have {
		return ErrUserNotRoleExist
	}

	delete(u, rid)
	return nil
}

// check user has permission
func (auth *Auth) Permit(uid, pid string) (bool, error) {
	auth.RLock()
	defer auth.RUnlock()

	u, have := auth.users[uid]
	if !have {
		return false, ErrUserNotExist
	}

	p, err := auth.GetPerm(pid)
	if err != nil {
		return false, err
	}

	isPerm := rolesPermit(u, p)
	return isPerm, nil
}

// check role has permission
func (auth *Auth) RolePermit(rid, pid string) (bool, error) {
	p, err := auth.GetPerm(pid)
	if err != nil {
		return false, err
	}

	r, err := auth.GetRole(rid)
	if err != nil {
		return false, err
	}

	u := map[string]Role{"role_check": r}
	return rolesPermit(u, p), nil
}

func (auth *Auth) GetAllUsers() map[string]User {
	return auth.users
}

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

func (auth *Auth) AddRole(uid, rid string) error {
	auth.Lock()
	defer auth.Unlock()

	u, have := auth.users[uid]
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

	u[rid] = r
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

/**
 *
 */

func (auth *Auth) SetParents(id string, parents []string) error {
	auth.Lock()
	defer auth.Unlock()
	if _, ok := auth.roles[id]; !ok {
		return ErrRoleNotExist
	}
	for _, parent := range parents {
		if _, ok := auth.roles[parent]; !ok {
			return ErrRoleNotExist
		}
	}
	if _, ok := auth.parents[id]; !ok {
		auth.parents[id] = make(map[string]struct{})
	}
	for _, parent := range parents {
		auth.parents[id][parent] = struct{}{}
	}
	return nil
}

func (auth *Auth) GetParents(id string) ([]string, error) {
	auth.Lock()
	defer auth.Unlock()
	if _, ok := auth.roles[id]; !ok {
		return nil, ErrRoleNotExist
	}
	ids, ok := auth.parents[id]
	if !ok {
		return nil, nil
	}
	var parents []string
	for parent := range ids {
		parents = append(parents, parent)
	}
	return parents, nil
}

func (auth *Auth) SetParent(id string, parent string) error {
	auth.Lock()
	defer auth.Unlock()
	if _, ok := auth.roles[id]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := auth.roles[parent]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := auth.parents[id]; !ok {
		auth.parents[id] = make(map[string]struct{})
	}
	var empty struct{}
	auth.parents[id][parent] = empty
	return nil
}

func (auth *Auth) RemoveParent(id string, parent string) error {
	auth.Lock()
	defer auth.Unlock()
	if _, ok := auth.roles[id]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := auth.roles[parent]; !ok {
		return ErrRoleNotExist
	}
	delete(auth.parents[id], parent)
	return nil
}

func (auth *Auth) Add(r Role) (err error) {
	auth.Lock()
	defer auth.Unlock()
	if _, ok := auth.roles[r.getId()]; !ok {
		auth.roles[r.getId()] = r
	} else {
		err = ErrRoleExist
	}
	return
}

func (auth *Auth) Remove(id string) (err error) {
	auth.Lock()
	defer auth.Unlock()
	if _, ok := auth.roles[id]; ok {
		delete(auth.roles, id)
		for rid, parents := range auth.parents {
			if rid == id {
				delete(auth.parents, rid)
				continue
			}
			for parent := range parents {
				if parent == id {
					delete(auth.parents[rid], id)
					break
				}
			}
		}
	} else {
		err = ErrRoleNotExist
	}
	return
}

func (auth *Auth) Get(id string) (r Role, parents []string, err error) {
	auth.RLock()
	defer auth.RUnlock()

	var ok bool
	if r, ok = auth.roles[id]; ok {
		for parent := range auth.parents[id] {
			parents = append(parents, parent)
		}
	} else {
		err = ErrRoleNotExist
	}
	return
}

func (auth *Auth) IsGranted(id string, p Permission, assert AssertionFunc) (rslt bool) {
	auth.RLock()
	rslt = auth.isGranted(id, p, assert)
	auth.RUnlock()
	return
}

func (auth *Auth) isGranted(id string, p Permission, assert AssertionFunc) bool {
	if assert != nil && !assert(auth, id, p) {
		return false
	}
	return auth.recursionCheck(id, p)
}

func (auth *Auth) recursionCheck(id string, p Permission) bool {
	if role, ok := auth.roles[id]; ok {
		if role.permit(p) {
			return true
		}
		if parents, ok := auth.parents[id]; ok {
			for pID := range parents {
				if _, ok := auth.roles[pID]; ok {
					if auth.recursionCheck(pID, p) {
						return true
					}
				}
			}
		}
	}
	return false
}

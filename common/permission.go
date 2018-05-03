package common

import (
	"strings"
)

type firstPermission struct {
	id      string `json:"id"`
	descrip string `json:"descrip"`
	sep     string `json:"sep"`
}

func NewFirstP(id, descrip string) *firstPermission {
	return &firstPermission{id, descrip, FirstSep}
}

func (p *firstPermission) getDes() string {
	return p.descrip
}

func (p *firstPermission) getId() string {
	return p.id
}

func (p *firstPermission) setDes(descrip string) {
	p.descrip = descrip
}

func (p *firstPermission) match(a Permission) bool {
	if p.id == a.getId() || p.descrip == a.getDes() {
		return true
	}
	q, ok := a.(*firstPermission)
	if !ok {
		return false
	}
	players := strings.Split(p.descrip, p.sep)
	qlayers := strings.Split(q.descrip, q.sep)

	// 可以包含
	if len(players) > len(qlayers) {
		return false
	}
	for k, pv := range players {
		if pv != qlayers[k] {
			return false
		}
	}
	return true
}

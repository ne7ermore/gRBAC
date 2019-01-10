package plugin

import (
	"time"
)

type PermissionPools interface {
	Gather() ([]Permission, error)
	New(string, string) (string, error)
	Get(string) (Permission, error)
	GetByDesc(string) (Permission, error)
	Update(string, map[string]string) error
	Gets(int, int, string) ([]Permission, error)
	Close()
	Counts() int
}

type Permission interface {
	Getid() string
	GetName() string
	GetDescrip() string
	GetSep() string
	GetCreateTime() time.Time
	GetUpdateTime() time.Time
}

package plugin

import (
	"time"
)

type RolePools interface {
	Gather() ([]Role, error)
	New(string) (string, error)
	Get(string) (Role, error)
	Update(string, map[string]string) error
	Gets(int, int, string) ([]Role, error)
	GetByName(string) (Role, error)
	Close()
	Counts() int
}

type Role interface {
	Getid() string
	GetName() string
	GetPermissions() string
	GetCreateTime() time.Time
	GetUpdateTime() time.Time
}

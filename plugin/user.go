package plugin

import (
	"time"
)

type UserPools interface {
	Gather() ([]User, error)
	New(string) (string, error)
	Get(string) (User, error)
	GetByUid(string) (User, error)
	Update(string, map[string]string) error
	Gets(int, int, string) ([]User, error)
	Close()
	Counts() int
}

type User interface {
	Getid() string
	GetUserId() string
	GetRoles() string
	GetCreateTime() time.Time
	GetUpdateTime() time.Time
}

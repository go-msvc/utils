package registry

import "time"

type Registry interface {
	Find(filters []string, limit int) (entries []RegistryEntry, nrFound int, err error)
	Register(id string, names []string) error
	Unregister(id string) error
}

type RegistryEntry struct {
	ID        string
	Names     []string
	StartTime time.Time
	LastTime  time.Time
}

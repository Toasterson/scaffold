package storage

import "github.com/satori/go.uuid"

type Storeable interface {
	SetUUID(u uuid.UUID) error
	HasUUID() bool
	GetUUID() uuid.UUID
}

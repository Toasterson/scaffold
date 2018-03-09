package storage

import "github.com/satori/go.uuid"

type Storer interface {
	Create(object Storeable) error

	All(results *[]interface{}, obj interface{}) error

	Has(id uuid.UUID) bool

	Read(id uuid.UUID, object interface{}) error

	Update(object Storeable) error

	Delete(id uuid.UUID) error
}

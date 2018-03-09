package storage

import (
	"fmt"

	"github.com/satori/go.uuid"
)

type MemStorer struct {
	store map[uuid.UUID]interface{}
}

func (m *MemStorer) Create(object Storeable) error {
	var key uuid.UUID
	if object.HasUUID() {
		key = object.GetUUID()
	} else {
		key = uuid.NewV4()
		object.SetUUID(key)
	}
	if _, ok := m.store[key]; ok {
		return fmt.Errorf("duplicate definition of uuid %s, this is probably malicious", key)
	}
	m.store[key] = object
	return nil
}

func (m *MemStorer) All(results *[]interface{}, obj interface{}) error {
	for _, value := range m.store {
		*results = append(*results, value)
	}
	return nil
}

func (m *MemStorer) Has(id uuid.UUID) bool {
	_, ok := m.store[id]
	return ok
}

func (m *MemStorer) Read(id uuid.UUID, object interface{}) error {
	if value, ok := m.store[id]; ok {
		object = value
		return nil
	} else {
		return fmt.Errorf("could not find %s", id)
	}
}

func (m *MemStorer) Update(object Storeable) error {
	m.store[object.GetUUID()] = object
	return nil
}

func (m *MemStorer) Delete(id uuid.UUID) error {
	if m.Has(id) {
		delete(m.store, id)
		return nil
	} else {
		return fmt.Errorf("%s does not exist", id)
	}
}

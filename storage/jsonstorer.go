package storage

import (
	"encoding/json"
	"reflect"

	"github.com/satori/go.uuid"
)

func NewJsonStorer(dir, collection string) (*JsonStorer, error) {
	db, err := scribble.New(dir, nil)
	if err != nil {
		return nil, err
	}
	return &JsonStorer{db: db, collection: collection}, nil
}

type JsonStorer struct {
	db         *scribble.Driver
	collection string
}

func (j *JsonStorer) Create(object Storeable) error {
	if err := j.db.Write(j.collection, object.GetUUID().String(), object); err != nil {
		return err
	}
	return nil
}

func (j *JsonStorer) All(results *[]interface{}, obj interface{}) error {
	records, err := j.db.ReadAll(j.collection)
	if err != nil {
		return err
	}
	for _, record := range records {
		unmarshalled := reflect.New(reflect.TypeOf(obj)).Interface()
		if err := json.Unmarshal([]byte(record), &unmarshalled); err != nil {
			return err
		}
		*results = append(*results, unmarshalled)
	}
	return nil
}

func (j *JsonStorer) Has(id uuid.UUID) bool {
	if err := j.db.Read(j.collection, id.String(), nil); err != nil {
		return false
	}
	return true
}

func (j *JsonStorer) Read(id uuid.UUID, object interface{}) error {
	if err := j.db.Read(j.collection, id.String(), object); err != nil {
		return err
	}
	return nil
}

func (j *JsonStorer) Update(object Storeable) error {
	if err := j.db.Write(j.collection, object.GetUUID().String(), object); err != nil {
		return err
	}
	return nil
}

func (j *JsonStorer) Delete(id uuid.UUID) error {
	if err := j.db.Delete(j.collection, id.String()); err != nil {
		return err
	}
	return nil
}

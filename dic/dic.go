package dic

import (
	"errors"
)

type Dictionary map[string]string

var errNotFound = errors.New("not found")

func (d Dictionary) Search(key string) (string, error) {
	value, exist := d[key]

	if exist {
		return value, nil
	}
	return "", errNotFound
}

func (d Dictionary) Add(key, value string) error {
	_, exist := d[key]
	if exist {
		return errors.New("already exist")
	}
	d[key] = value

	return nil
}

func (d Dictionary) Update(key, value string) {
	d[key] = value
}

func (d Dictionary) Delete(key string) {
	delete(d, key)
}
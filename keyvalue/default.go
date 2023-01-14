package keyvalue

import (
	"encoding/json"
	"reflect"

	"github.com/go-msvc/errors"
)

func init() {
	RegisterImplementation("default", MemConfig{})
}

type MemConfig struct{}

func (c MemConfig) Validate() error { return nil }

func (c MemConfig) Create() (Store, error) {
	return inMemStore{
		values: map[string]interface{}{},
	}, nil
}

type inMemStore struct {
	values map[string]interface{}
}

func (s inMemStore) Get(key string) (interface{}, error) {
	if v, ok := s.values[key]; ok {
		return v, nil
	}
	return nil, nil
}

func (s inMemStore) GetTmpl(key string, tmpl interface{}) (interface{}, error) {
	v, ok := s.values[key]
	if !ok {
		return nil, nil
	}
	jsonValue, err := json.Marshal(v)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot write value to JSON for conversion")
	}

	valueType := reflect.TypeOf(tmpl)
	defaultValue := reflect.ValueOf(tmpl)
	refCount := 0
	for valueType.Kind() == reflect.Ptr {
		refCount++
		valueType = valueType.Elem()
		defaultValue = defaultValue.Elem()
	}
	newPtrValue := reflect.New(valueType)
	newPtrValue.Elem().Set(defaultValue)
	if err := json.Unmarshal(jsonValue, newPtrValue.Interface()); err != nil {
		return nil, errors.Errorf("failed to parse JSON value")
	}
	if validator, ok := newPtrValue.Interface().(Validator); ok {
		if err := validator.Validate(); err != nil {
			return newPtrValue.Elem().Interface(), errors.Wrapf(err, "invalid value")
		}
	}
	for refCount > 0 {
		newPtrValue = newPtrValue.Addr()
		refCount--
	}
	return newPtrValue.Elem().Interface(), nil
} //GetTmpl()

func (s inMemStore) Set(key string, value interface{}) error {
	s.values[key] = value
	return nil
}

// for ms.UsedService interface
func (s inMemStore) Status() interface{} {
	return nil
}

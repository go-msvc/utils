package keyvalue

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/go-msvc/config"
	"github.com/go-msvc/errors"
)

func init() {
	config.RegisterConstructor("mem", MemConfig{})
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

func (s inMemStore) Del(key string) error {
	delete(s.values, key)
	return nil
}

func (s inMemStore) CtxGet(ctx context.Context, key string) (interface{}, error) {
	return s.Get(key)
}

func (s inMemStore) CtxGetTmpl(ctx context.Context, key string, tmpl interface{}) (interface{}, error) {
	return s.GetTmpl(key, tmpl)
}

func (s inMemStore) CtxSet(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return s.Set(key, value)
}

func (s inMemStore) CtxDel(ctx context.Context, key string) error {
	return s.Del(key)
}

// for ms.UsedService interface
func (s inMemStore) Status() interface{} {
	return nil
}

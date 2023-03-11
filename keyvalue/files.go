package keyvalue

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"time"

	"github.com/go-msvc/config"
	"github.com/go-msvc/errors"
)

func init() {
	config.RegisterConstructor("files", FilesConfig{})
}

type FilesConfig struct {
	Path string `json:"path" doc:"Directory where files are stored"`
}

func (c FilesConfig) Validate() error { return nil }

func (c FilesConfig) Create() (Store, error) {
	return filesStore{}, nil
}

type filesStore struct {
	path string
}

func (s filesStore) filename(key string) string {
	return s.path + "/" + key + ".json"
}

func (s filesStore) Get(key string) (interface{}, error) {
	fn := s.filename(key)
	f, err := os.Open(fn)
	if err != nil {
		return nil, errors.Errorf("cannot open file %s", fn)
	}
	defer f.Close()
	var value interface{}
	if err := json.NewDecoder(f).Decode(value); err != nil {
		return nil, errors.Errorf("failed to read JSON file %s", fn)
	}
	return value, nil
}

func (s filesStore) GetTmpl(key string, tmpl interface{}) (interface{}, error) {
	fn := s.filename(key)
	f, err := os.Open(fn)
	if err != nil {
		return nil, errors.Errorf("cannot open file %s", fn)
	}
	defer f.Close()
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
	if err := json.NewDecoder(f).Decode(newPtrValue.Interface()); err != nil {
		return nil, errors.Errorf("failed to read JSON file %s", fn)
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

func (s filesStore) Set(key string, value interface{}) error {
	fn := s.filename(key)
	f, err := os.Create(fn)
	if err != nil {
		return errors.Errorf("cannot create file %s", fn)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(value); err != nil {
		return errors.Errorf("failed encode JSON value")
	}
	return nil
}

func (s filesStore) Del(key string) error {
	fn := s.filename(key)
	err := os.Remove(fn)
	if err != nil && err != os.ErrNotExist {
		return errors.Errorf("cannot delete file %s", fn)
	}
	return nil
} //fileStore.Del()

func (s filesStore) CtxGet(ctx context.Context, key string) (interface{}, error) {
	return s.Get(key)
}

func (s filesStore) CtxGetTmpl(ctx context.Context, key string, tmpl interface{}) (interface{}, error) {
	return s.GetTmpl(key, tmpl)
}

func (s filesStore) CtxSet(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return s.Set(key, value)
}

func (s filesStore) CtxDel(ctx context.Context, key string) error {
	return s.Del(key)
}

// for ms.UsedService interface
func (s filesStore) Status() interface{} {
	return nil
}

type Validator interface {
	Validate() error
}

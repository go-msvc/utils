package keyvalue

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/go-msvc/errors"
)

func init() {
	RegisterImplementation("files", FilesConfig{})
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

// for ms.UsedService interface
func (s filesStore) Status() interface{} {
	return nil
}

type Validator interface {
	Validate() error
}

package keyvalue

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-msvc/errors"
	"github.com/go-msvc/utils/ms"
	"github.com/go-msvc/utils/stringutils"
)

type Constructor interface {
	ms.Config
	Create() (Store, error)
}

type Store interface {
	ms.UsedService
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
}

//todo: define standard errors
//	store not reachable
//	invalid key
//	key not found

var (
	//the registered implementations store a template of its config
	implementationNames  = []string{}
	implementationByName = map[string]Constructor{}
	implementationsMutex sync.Mutex
)

func RegisterImplementation(name string, tmplConfig Constructor) {
	implementationsMutex.Lock()
	defer implementationsMutex.Unlock()
	if !stringutils.IsSnakeCase(name) {
		panic(fmt.Sprintf("KeyValue implementation(%s) must be named with snake_case", name))
	}
	if _, ok := implementationByName[name]; ok {
		panic(fmt.Sprintf("duplicate implementation(%s)", name))
	}
	if tmplConfig == nil {
		panic(fmt.Sprintf("implementation(%s) == nil", name))
	}
	implementationByName[name] = tmplConfig
	implementationNames = append(implementationNames, name)
}

type Config map[string]interface{}

//include this in your service config to create the kv for you
//during service configuration
func (c *Config) Validate() error {
	if len(*c) != 1 {
		return errors.Errorf("%d items, expecting one", len(*c))
	}

	//get the configured implementation name and value
	var implementationName string
	var configValue interface{}
	for implementationName, configValue = range *c {
	}

	//get the registered implementation by name
	tmplConfig, ok := implementationByName[implementationName]
	if !ok {
		return errors.Errorf("unknown implementation(%s) expecting %s", implementationName, strings.Join(implementationNames, "|"))
	}

	//make a copy of the tmplConfig, then parse the value into that struct,
	//overwriting only fields that are defined
	configValuePtr := reflect.New(reflect.TypeOf(tmplConfig))
	configValuePtr.Elem().Set(reflect.ValueOf(tmplConfig))

	jsonConfigValue, _ := json.Marshal(configValue)
	if err := json.Unmarshal(jsonConfigValue, configValuePtr.Interface()); err != nil {
		return errors.Errorf("cannot decode JSON into %T", tmplConfig)
	}
	if validator, ok := configValuePtr.Interface().(ms.Validator); !ok {
		return errors.Errorf("%T does not implement Validator", tmplConfig)
	} else {
		if err := validator.Validate(); err != nil {
			return errors.Wrapf(err, "invalid %s", implementationName)
		}
	}

	constructor, ok := configValuePtr.Elem().Interface().(Constructor)
	if !ok {
		return errors.Errorf("%T does not implement KeyValueConfig", configValuePtr.Elem().Interface())
	}

	//replace value with validated config ready to create when all config was loaded
	(*c)[implementationName] = constructor
	return nil
} //Config.Validate()

func (c Config) Create() (Store, error) {
	var implementationName string
	var configValue interface{}
	for implementationName, configValue = range c {
		//empty loop to get the one and only name and value from the map
	}
	kvc, ok := configValue.(Constructor)
	if !ok {
		panic(fmt.Sprintf("%T is not KeyValueConfig", configValue))
	}
	store, err := kvc.Create()
	if err != nil {
		panic(fmt.Sprintf("%s failed to create: %+v", implementationName, err))
	}
	return store, nil
}

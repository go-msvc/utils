package ms

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/go-msvc/errors"
)

//Server is a synchronous micro-service interface
//it gets requests, call the appropriate handler then sends a response
//examples are:
//	HTTP REST server
//	NATS subscription on request topics (it publishes a response on a topid specified in the request)
//
//see also Consumer for async micro-services

//server may be sync or async
type ServerConfig interface {
	Create(MicroService) (Server, error)
}

type Server interface {
	Serve() error
}

type serverConfig map[string]interface{} //expected value ServerConfig

func (scObj *serverConfig) Validate() error {
	if len(registeredServerImplementations) == 0 {
		return errors.Errorf("no server implementations included - seems anonymous imports were not done")
	}
	log.Infof("Validating sc %+v", scObj)
	if len(*scObj) != 1 {
		names := []string{}
		for n := range *scObj {
			names = append(names, n)
		}
		return errors.Errorf("%d elements %v instead of 1", len(*scObj), names)
	}
	var configuredName string
	var configuredValue interface{}
	for configuredName, configuredValue = range *scObj {
		//empty loop to get the one and only name and value from the map
	}

	//name must be registered server implementation
	tmpl, ok := registeredServerImplementations[configuredName]
	if !ok {
		return errors.Errorf("config server.%s is not registered as a server implementation", configuredName)
	}

	t := reflect.TypeOf(tmpl)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	configPtrValue := reflect.New(t)

	jsonValue, _ := json.Marshal(configuredValue)
	if err := json.Unmarshal(jsonValue, configPtrValue.Interface()); err != nil {
		panic(fmt.Sprintf("failed to decode JSON configured value for server.%s into %v: %+v", configuredName, t, err))
	}
	if err := configPtrValue.Interface().(Validator).Validate(); err != nil {
		panic(fmt.Sprintf("invalid server.%s: %+v", configuredName, err))
	}

	//replace the generic interface{} value from config source
	//with the validated and structured ServerConfig value
	(*scObj)[configuredName] = configPtrValue.Interface().(ServerConfig)
	return nil
}

var (
	registeredServerImplementations = map[string]ServerConfig{}
	registeredServerMutex           sync.Mutex
)

func RegisteredServerImplementation(name string, tmpl ServerConfig) {
	registeredServerMutex.Lock()
	defer registeredServerMutex.Unlock()
	if _, ok := registeredServerImplementations[name]; ok {
		panic(fmt.Sprintf("duplicate server implementation registered for \"%s\"", name))
	}
	if tmpl == nil {
		panic(fmt.Sprintf("server implementation \"%s\" cannot be nil", name))
	}
	registeredServerImplementations[name] = tmpl
}

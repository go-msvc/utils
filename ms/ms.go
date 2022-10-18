package ms

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-msvc/errors"
	"github.com/go-msvc/utils/stringutils"
	"github.com/stewelarend/logger"
)

var log = logger.New().WithLevel(logger.LevelDebug)

type MicroService interface {
	Configure()

	//call one of the following in main() after Configure():
	//	Serve() if this is a synchronouse micro-service
	//	Consume() if this is an asynchronous micro-service
	Serve()
	Consume()

	Oper(name string) (Operation, bool)
	OperNames() []string
	NewContext() context.Context
}

type Option func(s *microService) error

func New(name string, options ...Option) microService {
	ms := microService{
		configs:   map[string]config{},
		opers:     map[string]oper{},
		operNames: []string{},
	}
	minimalOptions := []Option{
		WithConfig("server", nil, &serverConfig{}),
		withOper("_version", ms.operVersion),
		withOper("_doc", ms.operDoc),
		withOper("_status", ms.operStatus),
		withOper("_stop", ms.operStop),
	}
	for _, option := range minimalOptions {
		if err := option(&ms); err != nil {
			panic(fmt.Sprintf("failed to define micro-service: %+v", err))
		}
	}
	for _, option := range options {
		if err := option(&ms); err != nil {
			panic(fmt.Sprintf("failed to apply option: %+v", err))
		}
	}
	return ms
}

func WithConfig(name string, ctxID interface{}, tmpl Config) Option {
	return func(ms *microService) error {
		if !stringutils.IsSnakeCase(name) {
			panic(fmt.Sprintf("cannot add config name \"%s\" because it is not written as snake_case", name))
		}
		if _, ok := ms.configs[name]; ok {
			panic(fmt.Sprintf("duplicate config name \"%s\"", name))
		}
		ms.configs[name] = config{
			ctxID: ctxID,
			value: tmpl,
		}
		return nil
	}
}

//just config values
type Config interface {
	Validate() error
}

type ConfiguredConstructor interface {
	Config
	Create() (any, error)
}

type microService struct {
	configs           map[string]config
	opers             map[string]oper
	operNames         []string
	namedConfigValues map[string]interface{} //fetched from source in ms.Configure()
}

var (
	VersionString   = "dev"
	BuildTimeString = "undefined"
)

type VersionInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
}

func (ms microService) operVersion(ctx context.Context) (VersionInfo, error) {
	return VersionInfo{
		Version:   VersionString,
		BuildTime: BuildTimeString,
	}, nil
}
func (ms microService) operStatus(ctx context.Context) (interface{}, error) { panic("NYI") }

type DocRequest struct {
	Path string `json:"path"`
}

func (ms microService) operDoc(ctx context.Context, req DocRequest) (res string, err error) {
	//todo: output json
	panic("NYI")
}

func (ms microService) operStop(ctx context.Context) error { panic("NYI") }

func (ms *microService) Configure() {
	//fetch config (i.e. load from file for dev or fetch remote in clustered deployment)
	//this is all config - never changes and cannot add to it once this is done
	var err error
	ms.namedConfigValues, err = fetchConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to fetch config: %+v", err))
	}

	//get all registered configurations
	for n, config := range ms.configs {
		//see if this name has a configured value
		if v, ok := ms.namedConfigValues[n]; ok {
			log.Debugf("\"%s\" is configured: %+v", n, v)
			//got a named config value
			//marshal to json than unmarshal over the tmpl for this name
			jsonValue, err := json.Marshal(v)
			if err != nil {
				panic(fmt.Sprintf("failed to JSON encode configured value for %s:%+v: %+v", n, v, err))
			}

			//create new copy of tmpl to parse into
			t := reflect.TypeOf(config.value)
			for t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			configPtrValue := reflect.New(t)
			if err := json.Unmarshal(jsonValue, configPtrValue.Interface()); err != nil {
				log.Errorf("JSON: %s", string(jsonValue))
				panic(fmt.Sprintf("failed to decode JSON configured value for \"%s\":%s into %v: %+v", n, string(jsonValue), t, err))
			}
			if err := configPtrValue.Interface().(Validator).Validate(); err != nil {
				panic(fmt.Sprintf("invalid \"%s\":%s: %+v", n, string(jsonValue), err))
			}
			config.value = configPtrValue.Interface().(Config)
			ms.configs[n] = config
		} else {
			log.Debugf("\"%s\" is NOT configured", n)
		}

		//log config
		{
			defaultConfig, _ := json.Marshal(config.value)
			log.Infof("config \"%s\":%s", n, &defaultConfig)
		}

		if constructor, ok := config.value.(ConfiguredConstructor); ok {
			config.created, err = constructor.Create()
			if err != nil {
				panic(fmt.Sprintf("failed to construct %s from config: %+v", n, err))
			}
			log.Infof("created configured \"%s\": (%T)", n, config.created)
			ms.configs[n] = config
		} else {
			log.Infof("config \"%s\" is just a value (not constructor)", n)
		}
	}
} //microService.Configure()

func (ms microService) Serve() {
	sc := ms.configs["server"].value.(*serverConfig)
	log.Infof("GOT SERVER CONFIG: %T: %+v", sc, sc)
	var serverName string
	var serverConfig ServerConfig
	for n, v := range *sc {
		serverName = n
		serverConfig = v.(ServerConfig)
	}
	log.Infof("  creating server.%s: (%T)", serverName, serverConfig)

	//create the configured server
	// serverImplementationName, configValue, err := ms.GetNamedConfig("server")
	// if err != nil {
	// 	panic(fmt.Sprintf("cannot get config \"server\": %+v", err))
	// }
	// sc, ok := configValue.(ServerConfig)
	// if !ok {
	// 	panic(fmt.Sprintf("server config value type %T != ServerConfig", configValue))
	// }
	s, err := serverConfig.Create(&ms)
	if err != nil {
		panic(fmt.Sprintf("failed to create server \"%s\": %+v", serverName, err))
	}

	//run the server
	//todo: for now simple HTTP server
	log.Infof("Starting server \"%s\"", serverName)
	if err := s.Serve(); err != nil {
		panic(fmt.Sprintf("server \"%s\" failed: %+v", serverName, err))
	}
	log.Infof("server \"%s\" terminated", serverName)
} //microService.Serve()

func (ms microService) Consume() {
	panic("NYI")
}

func (ms microService) GetNamedConfig(name string) (implName string, implConfigValue interface{}, err error) {
	if v, ok := ms.namedConfigValues[name]; ok {
		//got the config, expecting one named value to select an implementation
		obj, ok := v.(map[string]interface{})
		if !ok {
			return "", nil, errors.Errorf("configured value for \"%s\" is not a JSON object", name)
		}
		if len(obj) != 1 {
			return "", nil, errors.Errorf("configured value for \"%s\" has %d instead of 1 items", name, len(obj))
		}
		for implName, implConfigValue = range obj {
			//empty loop to get the one and only name and value from the map
		}
		log.Debugf("using config \"%s\":{\"%s\":(%T)%+v}", name, implName, implConfigValue, implConfigValue)
		return implName, implConfigValue, nil
	}
	return "", nil, errors.Errorf("\"%s\" is not configured", name)
} //microService.GetNamesConfig()

func (ms microService) Oper(name string) (Operation, bool) {
	if o, ok := ms.opers[name]; ok {
		return o, true
	}
	return nil, false
}

func (ms microService) OperNames() []string {
	return ms.operNames
}

func (ms microService) NewContext() context.Context {
	ctx := context.Background()
	for _, config := range ms.configs {
		if config.ctxID != nil {
			if config.created != nil {
				ctx = context.WithValue(ctx, config.ctxID, config.created)
			} else {
				ctx = context.WithValue(ctx, config.ctxID, config.value)
			}
		}
	}
	return ctx
}

var errorInterfaceType = reflect.TypeOf((*error)(nil)).Elem()
var contextInterfaceType = reflect.TypeOf((*context.Context)(nil)).Elem()

type Validator interface {
	Validate() error
}

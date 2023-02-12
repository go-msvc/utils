package ms

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-msvc/config"
	"github.com/go-msvc/errors"
	"github.com/go-msvc/logger"
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

func New(options ...Option) microService {
	ms := microService{
		opers:     map[string]oper{},
		operNames: []string{},
	}
	minimalOptions := []Option{
		//get this from github.com/go-msvc/Config with server construction WithConfig("server", nil, &serverConfig{}),
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
	config.MustConstruct("ms.server", reflect.TypeOf((*Server)(nil)).Elem())
	return ms
}

type microService struct {
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
	if err := config.Load(); err != nil {
		panic(fmt.Sprintf("failed to load config: %+v", err))
	}

	//todo: capture config used, make it available in local file to reload if
	//config source is broken next time?
	//todo: generate docs of values and schema
} //microService.Configure()

func (ms *microService) Serve() {
	s := config.Get("ms.server").(Server)

	//run the server
	log.Infof("Starting server ...")
	if err := s.Serve(ms); err != nil {
		panic(fmt.Sprintf("server failed: %+v", err))
	}
	log.Infof("server terminated")
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
	return ctx
}

var errorInterfaceType = reflect.TypeOf((*error)(nil)).Elem()
var contextInterfaceType = reflect.TypeOf((*context.Context)(nil)).Elem()

type Validator interface {
	Validate() error
}

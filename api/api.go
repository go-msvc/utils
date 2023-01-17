package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-msvc/config"
	"github.com/go-msvc/errors"
	"github.com/go-msvc/logger"
	"github.com/go-msvc/utils/ms"
)

var log = logger.New().WithLevel(logger.LevelDebug)

func init() {
	config.RegisterConstructor("http", ServerConfig{})
}

type ServerConfig struct {
	Addr string `json:"addr"`
}

func (c ServerConfig) Validate() error {
	if c.Addr == "" {
		return errors.Errorf("missing addr")
	}
	return nil
}

func (c ServerConfig) Create() (ms.Server, error) {
	return &server{
		config: c,
		svc:    nil, //defined in Serve()
	}, nil
}

type server struct {
	config ServerConfig
	svc    ms.MicroService
}

func (s *server) Serve(svc ms.MicroService) error {
	s.svc = svc
	if err := http.ListenAndServe(s.config.Addr, s); err != nil {
		return errors.Wrapf(err, "HTTP server failed")
	}
	return nil
}

func (s server) ServeHTTP(httpRes http.ResponseWriter, httpReq *http.Request) {
	log.Debugf("HTTP %s %s", httpReq.Method, httpReq.URL.Path)

	//todo: throttle source

	//todo: auth

	//todo: log

	if httpReq.Method != http.MethodPost {
		http.Error(httpRes, "expecting POST for all operations", http.StatusMethodNotAllowed)
		return
	}

	//identify the operation
	operName := httpReq.URL.Path[1:] //skip leading '/'
	oper, ok := s.svc.Oper(operName)
	if !ok {
		http.Error(httpRes, fmt.Sprintf("unknown operation \"%s\"", operName), http.StatusNotFound)
		return
	}

	var reqPtrValue reflect.Value
	if oper.ReqType() != nil {
		reqPtrValue = reflect.New(oper.ReqType())
		ct := httpReq.Header.Get("Content-Type")
		if ct == "" {
			http.Error(httpRes, fmt.Sprintf("missing header Content-Type"), http.StatusBadRequest)
			return
		}
		switch ct {
		case "application/json":
			if err := json.NewDecoder(httpReq.Body).Decode(reqPtrValue.Interface()); err != nil {
				http.Error(httpRes, fmt.Sprintf("invalid JSON in request body: %+v", err), http.StatusBadRequest)
				return
			}
		default:
			http.Error(httpRes, fmt.Sprintf("Content-Type:%s not supported", ct), http.StatusBadRequest)
			return
		}
		if validator, ok := reqPtrValue.Interface().(ms.Validator); ok {
			if err := validator.Validate(); err != nil {
				http.Error(httpRes, fmt.Sprintf("invalid %s request: %+v", operName, err), http.StatusBadRequest)
				return

			}
		}
	}

	ctx := context.Background()
	res, err := oper.Handle(ctx, reqPtrValue.Elem().Interface())
	if err != nil {
		http.Error(httpRes, fmt.Sprintf("%s failed: %+v", operName, err), http.StatusInternalServerError)
		return
	}

	if res != nil {
		log.Debugf("%s(%+v) -> %+v", operName, reqPtrValue.Elem().Interface(), res)
		jsonRes, err := json.Marshal(res)
		if err != nil {
			http.Error(httpRes, fmt.Sprintf("%s response encoding failed: %+v", operName, err), http.StatusInternalServerError)
			return
		}
		httpRes.Header().Set("Content-Type", "application/json")
		httpRes.Write(jsonRes)
	}
	return
}

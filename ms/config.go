package ms

import (
	"encoding/json"
	"os"

	"github.com/go-msvc/errors"
)

type config struct {
	ctxID   interface{}
	value   Config
	created interface{} //output of optional Create()
}

//fetch config in ms.Configure()
//called only once
func fetchConfig() (map[string]interface{}, error) {
	source := os.Getenv("MS_CONFIG_SOURCE")
	if source == "" {
		source = "./conf/config.json"
	}

	//todo: configurable source, allow etcd or http get etc...

	//todo: env for config version, else latest

	f, err := os.Open(source)
	if err != nil {
		log.Infof("Using all default configs")
		return nil, nil
		//return nil, errors.Wrapf(err, "failed to open config file \"%s\"", source)
	}
	defer f.Close()

	var configValue map[string]interface{}
	if err := json.NewDecoder(f).Decode(&configValue); err != nil {
		return nil, errors.Wrapf(err, "failed to decode config as JSON object")
	}
	return configValue, err
}

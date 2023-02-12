package ms

import (
	"context"
	"time"

	"github.com/go-msvc/errors"
)

type CtxID struct{}

type Address struct {
	Domain    string `json:"domain,omitempty" doc:"Domain name e.g. api.mystuff.co.za"`
	Operation string `json:"operation,omitempty" doc:"Operation name within the domain"`
}

func (a Address) Validate() error {
	if a.Domain == "" {
		return errors.Errorf("missing domain")
	}
	if a.Operation == "" {
		return errors.Errorf("missing operation")
	}
	return nil
} //Address.Validate()

type Client interface {
	Sync(
		ctx context.Context,
		serviceAddress Address,
		ttl time.Duration,
		req interface{},
		resTmpl interface{},
	) (
		res interface{},
		err error,
	)

	ASync(addr Address, req interface{}) (err error)
	Send(addr Address, ttl time.Duration, req interface{}) (err error)
}

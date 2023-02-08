package ms

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-msvc/utils/stringutils"
)

type Operation interface {
	Name() string
	ReqType() reflect.Type
	Handle(ctx context.Context, reqData interface{}) (resData interface{}, err error)
}

// oper implements Operation interface
type oper struct {
	name      string
	funcValue reflect.Value
	reqType   reflect.Type
}

// WithOper() enforces names to be snake_case
// while internally we call withOper() to allow leading underscore, not allowed for service's own operations
func WithOper(name string, handlerFunc interface{}) Option {
	if !stringutils.IsSnakeCase(name) {
		panic(fmt.Sprintf("cannot add operation name \"%s\" because it is not written as snake_case", name))
	}
	return withOper(name, handlerFunc)
}

func withOper(name string, handlerFunc interface{}) Option {
	return func(ms *microService) error {
		//todo: restrict name also to fit convention
		//todo: and use it always as part of the domain may be? e.g. name="jdbc" -> domain "ms-jdbc-someinst"
		if _, ok := ms.opers[name]; ok {
			panic(fmt.Sprintf("duplicate operation name \"%s\"", name))
		}
		o := oper{
			name:      name,
			funcValue: reflect.ValueOf(handlerFunc),
		}
		ft := reflect.TypeOf(handlerFunc)
		//first arg must be context
		if ft.NumIn() < 1 || !ft.In(0).Implements(contextInterfaceType) {
			panic(fmt.Sprintf("operation \"%s\" handler must take first argument of context.Context: %v != %v", name, ft.In(0), reflect.TypeOf(context.Background())))
		}
		//optional 2nd argument must be a request struct
		if ft.NumIn() > 1 {
			if ft.In(1).Kind() != reflect.Struct {
				panic(fmt.Sprintf("operation \"%s\" handler 2nd arg must be a struct type for the request", name))
			}
			o.reqType = ft.In(1)
		}
		//no more than 2 arguments allowed
		if ft.NumIn() > 2 {
			panic(fmt.Sprintf("operation \"%s\" handler may only take 2 arguments(context.Context, RequestStructType)", name))
		}
		//last result must be error
		if ft.NumOut() < 1 || !ft.Out(ft.NumOut()-1).Implements(errorInterfaceType) {
			panic(fmt.Sprintf("operation \"%s\" handler must return error as last result value", name))
		}
		//optional 2nd last result must be response - but could be any type
		//if ft.NumOut() > 1 {}
		//no more than 2 results expected
		if ft.NumOut() > 2 {
			panic(fmt.Sprintf("operation \"%s\" handler may not return more than (response, error)", name))
		}
		ms.opers[name] = o
		ms.operNames = append(ms.operNames, name)
		return nil
	}
} //withOper()

func (o oper) Name() string { return o.name }

func (o oper) ReqType() reflect.Type { return o.reqType }

func (o oper) Handle(ctx context.Context, req interface{}) (res interface{}, err error) {
	defer func() {
		if err != nil {
			log.Infof("%s.Handle (%T)%+v -> %+v", o.name, req, req, err)
		} else {
			log.Infof("%s.Handle (%T)%+v -> (%T)%+v", o.name, req, req, res, res)
		}
	}()

	args := []reflect.Value{
		reflect.ValueOf(ctx),
	}
	if o.reqType != nil {
		args = append(args, reflect.ValueOf(req))
	}
	results := o.funcValue.Call(args)

	//last result is error
	if results[len(results)-1].Interface() != nil {
		err = results[len(results)-1].Interface().(error)
		return
	}

	//res is optional, only if func returns more than just error, then first result
	if len(results) > 1 {
		res = results[0].Interface()
	}
	return
}

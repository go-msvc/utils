package ms

import "github.com/go-msvc/errors"

// Request is processed by an operation
type Request struct {
	//auth-ed user/role?
	//	Header map[string]interface{}
	Data interface{}
}

type Response struct {
	//	Header map[string]interface{} `json:"header,omitempty" doc:"Can be used in any context related to the operations"`
	Errors Errors      `json:"errors,omitempty" doc:"A stack of errors, [0] = final, [N-1] = original, none==success"`
	Data   interface{} `json:"data,omitempty" doc:"Operation specific request data"`
}

type Error struct {
	Code    string `json:"code" doc:"Identifies the result and should be a consistent string meaning something specific, e.g. BAD_REQUEST or FileNotFound, using any format, but ideally no spaces and only alpha-numeric characters."`
	Details string `json:"details,omitempty" doc:"Optional details apart from the code, typically the error message"`
	Source  string `json:"source,omitempty" doc:"indicate the source that reported this, and could be package:file:line or any text to identify the source"`
}

type Errors []Error

//convert received errors to a golang error
func NewError(ee ...Error) error {
	var err error
	for i := len(ee) - 1; i >= 0; i-- {
		e := ee[i]
		if err == nil {
			err = errors.Errorf("%s: %s", e.Code, e.Details) //todo: do not include message.go as source
		} else {
			err = errors.Wrapf(err, "%s: %s", e.Code, e.Details) //todo: do not include message.go as source
		}
	}
	return err
}

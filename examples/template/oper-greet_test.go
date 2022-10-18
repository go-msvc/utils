package main

import (
	"context"
	"encoding/json"
	"testing"
)

type test struct {
	req GreetRequest
	exp *GreetResponse
	err error
}

//todo: add req validation checks

//todo: add cases where handler returns an error with/without a response

//todo: insist on http status code values, with custom values allowed in the correct ranges

//todo: provide standard session storage for stateless execution

//todo: provide standard auth (=OAuth2) interface with user, session, claim, permissions, ...

//todo: define called services as part of ms definition,
//	and integrate with called service's interface package
//	and put own messages in a package that can be imported
//	and export it to a json schema file which can be imported
//	into other tools and code generators

//todo: standard limiter
//todo: standard async queues and consumers
//todo: standard scheduler

//todo: test cases from json files including req/res and all external req/res
//	then draw seq diagram for that

func Test1(t *testing.T) {
	tests := []test{
		//example test cases
		//all req must be valid, as validation is checked before handler is called
		{req: GreetRequest{Name: ""}, exp: &GreetResponse{Message: "Hello "}, err: nil},
		{req: GreetRequest{Name: "Koos"}, exp: &GreetResponse{Message: "Hello Koos"}, err: nil},
		{req: GreetRequest{Name: "Piet"}, exp: &GreetResponse{Message: "Hello Piet"}, err: nil},
	}
	for _, test := range tests {
		ctx := context.Background()
		res, err := operGreet(ctx, test.req)
		if err != nil {
			t.Fatalf("test failed: %+v", err)
		}
		resJSON, _ := json.Marshal(res)
		expJSON, _ := json.Marshal(test.exp)
		if string(expJSON) != string(resJSON) {
			t.Fatalf("res: %+v != %+v", string(resJSON), string(expJSON))
		}
	}
}

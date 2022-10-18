package main

import (
	"context"
	"fmt"

	"github.com/go-msvc/errors"
)

type GreetRequest struct {
	Name string `json:"name"`
}

func (req GreetRequest) Validate() error {
	if req.Name == "" {
		return errors.Errorf("missing name")
	}
	return nil
}

type GreetResponse struct {
	Message string `json:"message"`
}

func operGreet(ctx context.Context, req GreetRequest) (*GreetResponse, error) {
	return &GreetResponse{
		Message: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}

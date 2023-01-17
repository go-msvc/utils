package main

import "context"

type AddGroupRequest struct{}

type AddGroupResponse struct{}

func addGroup(ctx context.Context, req AddGroupRequest) (*AddGroupResponse, error) {
	return nil, MyError{Message: "NYI"}
}

type GetGroupRequest struct{}

type GetGroupResponse struct{}

func getGroup(ctx context.Context, req AddGroupRequest) (*AddGroupResponse, error) {
	return nil, MyError{Message: "NYI"}
}

type UpdGroupRequest struct{}

type UpdGroupResponse struct{}

func updGroup(ctx context.Context, req AddGroupRequest) (*AddGroupResponse, error) {
	return nil, MyError{Message: "NYI"}
}

type DelGroupRequest struct{}

type DelGroupResponse struct{}

func delGroup(ctx context.Context, req AddGroupRequest) (*AddGroupResponse, error) {
	return nil, MyError{Message: "NYI"}
}

type FindGroupRequest struct{}

type FindGroupResponse struct{}

func findGroup(ctx context.Context, req AddGroupRequest) (*AddGroupResponse, error) {
	return nil, MyError{Message: "NYI"}
}

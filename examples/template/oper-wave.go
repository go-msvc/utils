package main

import "context"

type WaveRequest struct{}

type WaveResponse struct{}

func operWave(ctx context.Context, req WaveRequest) (*WaveResponse, MyError) {
	return nil, MyError{Message: "NYI"}
}

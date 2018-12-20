package testservice

import (
	"context"
	"strings"
	_ "testing"
	"time"
)

type underlyingSetter struct{}

func (s *underlyingSetter) Set(ctx context.Context, req *ReqSet) (*RespSet, error) {
	return &RespSet{Upper: strings.ToUpper(req.Str), Mod7: req.Num % 7}, nil
}

type underlyingWaiter struct{}

func (w *underlyingWaiter) Wait(ctx context.Context, req *ReqWait) (*Empty, error) {
	select {
	case <-ctx.Done():
		return &Empty{}, ctx.Err()
	case <-time.After(time.Duration(req.Duration)):
		return &Empty{}, nil
	}
}

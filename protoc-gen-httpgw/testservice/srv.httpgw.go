// Code generated by protoc-gen-httpgw
// source: srv.proto
// DO NOT EDIT!

/*
Package testservice is a http proxy.
*/

package testservice

import (
	"context"

	"github.com/pressly/chi"
	"github.com/qiwitech/graceful"
)

func NewSrvSetterHandler(srv SrvSetterInterface, c graceful.Codec) graceful.Handlerer {
	return AddSrvSetterHandlers(nil, srv, c)
}
func AddSrvSetterHandlers(mux graceful.Handlerer, srv SrvSetterInterface, c graceful.Codec) graceful.Handlerer {
	if mux == nil {
		mux = chi.NewMux()
	}

	mux.Handle("/Set", graceful.NewHandler(
		c,
		func() interface{} { return &ReqSet{} },
		func(ctx context.Context, args interface{}) (interface{}, error) { return srv.Set(ctx, args.(*ReqSet)) }))

	return mux
}

type SrvSetterHTTPClient struct {
	*graceful.Client
}

func NewSrvSetterHTTPClient(cl *graceful.Client) SrvSetterHTTPClient {
	return SrvSetterHTTPClient{
		Client: cl,
	}
}

func (cl SrvSetterHTTPClient) Set(ctx context.Context, args *ReqSet) (*RespSet, error) {
	var resp RespSet
	err := cl.Client.Call(ctx, "Set", args, &resp)
	return &resp, err
}

type SrvSetterInterface interface {
	Set(context.Context, *ReqSet) (*RespSet, error)
}

func NewSrvWaiterHandler(srv SrvWaiterInterface, c graceful.Codec) graceful.Handlerer {
	return AddSrvWaiterHandlers(nil, srv, c)
}
func AddSrvWaiterHandlers(mux graceful.Handlerer, srv SrvWaiterInterface, c graceful.Codec) graceful.Handlerer {
	if mux == nil {
		mux = chi.NewMux()
	}

	mux.Handle("/Wait", graceful.NewHandler(
		c,
		func() interface{} { return &ReqWait{} },
		func(ctx context.Context, args interface{}) (interface{}, error) {
			return srv.Wait(ctx, args.(*ReqWait))
		}))

	return mux
}

type SrvWaiterHTTPClient struct {
	*graceful.Client
}

func NewSrvWaiterHTTPClient(cl *graceful.Client) SrvWaiterHTTPClient {
	return SrvWaiterHTTPClient{
		Client: cl,
	}
}

func (cl SrvWaiterHTTPClient) Wait(ctx context.Context, args *ReqWait) (*Empty, error) {
	var resp Empty
	err := cl.Client.Call(ctx, "Wait", args, &resp)
	return &resp, err
}

type SrvWaiterInterface interface {
	Wait(context.Context, *ReqWait) (*Empty, error)
}

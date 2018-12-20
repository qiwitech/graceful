package integration

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/qiwitech/graceful"
	api "github.com/qiwitech/graceful/systest/api_v1"
	"github.com/qiwitech/graceful/systest/pluto"
)

func APIServiceChain(implProc pluto.ProcessingInterface, implStor pluto.StorageInterface) (addr string, stop func(context.Context) error, err error) {
	// listen
	var l net.Listener
	l, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return addr, stop, err
	}
	addr = "http://" + l.Addr().String()

	// Start HTTP-server
	httpsrv := graceful.NewServer()
	httpsrv.UseShutdownMiddleware(http.StatusServiceUnavailable, []byte("Service Status: unavailable!"))
	go httpsrv.Serve(l)
	// func to stop it
	stop = func(ctx context.Context) error {
		return httpsrv.Shutdown(ctx, time.Second)
	}

	// Processing client
	httpsrv.Mount("/proc", pluto.NewProcessingHandler(implProc, &graceful.ProtobufCodec{}))
	c, _ := graceful.NewClient(addr+"/proc",
		&graceful.ProtobufCodec{},
		breaker.New(3, 1, time.Millisecond))
	procClent := pluto.NewProcessingHTTPClient(c)

	// Storage client
	httpsrv.Mount("/stor", pluto.NewStorageHandler(implStor, &graceful.ProtobufCodec{}))
	c, _ = graceful.NewClient(addr+"/stor",
		&graceful.ProtobufCodec{},
		breaker.New(3, 1, time.Millisecond))
	storClent := pluto.NewStorageHTTPClient(c)

	// Combine Processing and Storage clients into API
	translator := NewAPItoPlutoTranslator(procClent, storClent)
	httpsrv.Mount("/api", api.NewAPIHandler(translator, &graceful.JSONCodec{}))
	addr += "/api"
	return addr, stop, err
}

func FullServiceChain(implProc pluto.ProcessingInterface, implStor pluto.StorageInterface) (client api.APIInterface, stop func(context.Context) error, err error) {
	var addr string
	// start http-server
	addr, stop, err = APIServiceChain(implProc, implStor)
	if err != nil {
		return client, stop, err
	}
	// API client
	c, _ := graceful.NewClient(addr,
		&graceful.JSONCodec{},
		breaker.New(3, 1, time.Millisecond))
	client = api.NewAPIHTTPClient(c)

	return client, stop, err
}

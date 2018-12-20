package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/qiwitech/graceful"
	"github.com/qiwitech/graceful/protoc-gen-httpgw/testservice"
)

var (
	fTimeout = flag.Duration("timeout", time.Second, "timeout")
	fListen  = flag.String("listen", "", "listen to addr")
	fConnect = flag.String("addr", "", "connect to addr")
	fCmd     = flag.String("cmd", "", "command: one of set|wait")
	fRepeat  = flag.Int("repeat", 1, "repeat request N times")
	fWait    = flag.Duration("wait", time.Second/10, "repeat interval")
)

var shutdown func()

func main() {
	flag.Parse()

	if *fListen != "" {
		listen()
		return
	}

	br := breaker.New(3, 1, 3*time.Second)

	gcl, err := graceful.NewClient(*fConnect, &graceful.ProtobufCodec{}, br)
	if err != nil {
		panic(err)
	}

	for i := 0; i < *fRepeat; i++ {
		if i != 0 {
			time.Sleep(*fWait)
		}

		ctx := context.TODO()
		var cancel func()
		if *fTimeout != 0 {
			ctx, cancel = context.WithTimeout(ctx, *fTimeout)
		}

		var resp interface{}

		switch *fCmd {
		case "set":
			cl := testservice.NewSrvSetterHTTPClient(gcl)
			resp, err = cl.Set(ctx, &testservice.ReqSet{Str: "str11qq", Num: 440})
		case "wait":
			cl := testservice.NewSrvWaiterHTTPClient(gcl)
			resp, err = cl.Wait(ctx, &testservice.ReqWait{Duration: time.Second.Nanoseconds()})
		default:
			panic("no such command")
		}
		if cancel != nil {
			cancel()
		}

		fmt.Printf("result: (%T) %v error %v\n", resp, resp, err)
	}
}

func listen() {
	l, err := net.Listen("tcp", *fListen)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Start server on %v\n", l.Addr())

	s := &Service{}

	srv := graceful.NewServer()
	srv.UseShutdownMiddleware(http.StatusServiceUnavailable, []byte("server is going to shutdown\n"))
	shutdown = func() {
		to := 5 * time.Second
		fmt.Printf("Going to shutdown (wait for %v)...\n", to)
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()
		srv.Shutdown(ctx, to)
	}

	mux := testservice.AddSrvSetterHandlers(nil, s, &graceful.ProtobufCodec{})
	mux = testservice.AddSrvWaiterHandlers(mux, s, &graceful.ProtobufCodec{})

	srv.Mount("/", mux)

	go shutdownHandler()

	if err = srv.Serve(l); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func shutdownHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	<-c
	shutdown()
}

type Service struct {
}

func (s *Service) Set(ctx context.Context, req *testservice.ReqSet) (*testservice.RespSet, error) {
	log.Printf("service Set: %+v", req)
	return &testservice.RespSet{Upper: strings.ToUpper(req.Str), Mod7: req.Num % 7}, nil
}

func (s *Service) Wait(ctx context.Context, req *testservice.ReqWait) (*testservice.Empty, error) {
	log.Printf("service Wait: %+v", req)
	time.Sleep(time.Duration(req.Duration))
	return nil, nil
}

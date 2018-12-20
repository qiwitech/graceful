package integration

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/qiwitech/graceful"
	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	var (
		testErr               = errors.New("All right")
		timeout time.Duration = time.Second
	)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error(err)
	}
	addr := "http://" + l.Addr().String()

	s := graceful.NewServer()
	go s.Serve(l)
	defer s.Stop()

	s.Handle("/slow", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(timeout * 4)
		graceful.Error(w, testErr, http.StatusInternalServerError)
	}))
	s.Handle("/quick", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(timeout / 4)
		graceful.Error(w, testErr, http.StatusInternalServerError)
	}))

	c, err := graceful.NewClient(addr, &graceful.JSONCodec{}, nil)
	if err != nil {
		t.Error(err)
	}

	// 1
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	err = c.Call(ctx, "slow", nil, nil)
	err = graceful.ErrorFromURL(err)
	cancel()
	assert.Equal(t, context.DeadlineExceeded, err)

	// 2
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	err = c.Call(ctx, "quick", nil, nil)
	err = graceful.ErrorFromURL(err)
	cancel()
	assert.Equal(t, testErr, err)

}

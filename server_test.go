package graceful

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pressly/chi"
	"github.com/stretchr/testify/assert"
)

func TestShutdown(t *testing.T) {
	shutdownMsg := "Service Status: unavailable!"
	simpleErr := errors.New("Test Error!")
	period := time.Millisecond

	r := chi.NewRouter()
	r.Get("/func", func(w http.ResponseWriter, req *http.Request) {
		Error(w, simpleErr, http.StatusInternalServerError)
	})

	srv := NewServer()
	srv.UseShutdownMiddleware(http.StatusServiceUnavailable, []byte(shutdownMsg))
	srv.Mount("/test", r)

	fl := NewFakeListener()

	serveErrCh := make(chan error, 2)
	go func() {
		serveErrCh <- srv.Serve(fl)
	}()

	cl, err := NewClient("/test", &JSONCodec{}, nil)
	assert.NoError(t, err)

	http.DefaultTransport = FakeRoudTripper{func(req *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		return w.Result(), nil
	}}

	err = cl.Call(context.TODO(), "func", nil, nil)
	assert.Equal(t, simpleErr, err)

	go func() {
		err1 := srv.Shutdown(context.TODO(), period)
		assert.Equal(t, http.ErrServerClosed, err1)
		serveErrCh <- err1
	}()
	begin := time.Now()
	for time.Since(begin) < period {
		time.Sleep(period / 4)
		err = cl.Call(context.TODO(), "func", nil, nil)
		assert.EqualError(t, err, shutdownMsg)
	}

	err1 := <-serveErrCh
	err2 := <-serveErrCh

	assert.Equal(t, err1, err2)
}

func TestTimeout(t *testing.T) {
	var (
		testErr               = errors.New("All right")
		timeout time.Duration = time.Second / 10
	)

	l := NewFakeListener()
	http.DefaultTransport = l.DefaultTransport

	s := NewServer()
	go s.Serve(l)
	defer s.Stop()

	s.Handle("/slow", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(timeout * 2)
		Error(w, testErr, http.StatusInternalServerError)
	}))
	s.Handle("/quick", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(timeout / 2)
		Error(w, testErr, http.StatusInternalServerError)
	}))

	c, err := NewClient("http://localhost", &JSONCodec{}, nil)
	if err != nil {
		t.Error(err)
	}

	// 1
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	err = c.Call(ctx, "slow", nil, nil)
	err = ErrorFromURL(err)
	cancel()
	assert.Equal(t, context.DeadlineExceeded, err)

	// 2
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	err = c.Call(ctx, "quick", nil, nil)
	err = ErrorFromURL(err)
	cancel()
	assert.Equal(t, testErr, err)

}

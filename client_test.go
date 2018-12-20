package graceful

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/gogo/protobuf/test"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	var (
		testData = test.NidOptNative{
			Field1:  -123,
			Field14: "dsfasdfasdf",
		}
		testErr = errors.New("Test Error!")
	)

	srv := NewServer()
	srv.Handle("/ok", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		data, err := json.Marshal(&testData)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))
	srv.Handle("/err", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		Error(w, testErr, http.StatusInternalServerError)
	}))

	cl, err := NewClient("/", &JSONCodec{}, nil)
	assert.NoError(t, err)
	cl.SetRequestProcessor(func(req *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		return w.Result(), nil
	})

	var resp test.NidOptNative

	err = cl.Call(context.TODO(), "ok", &test.NidOptNative{}, &resp)
	assert.NoError(t, err)
	assert.Equal(t, testData, resp)

	err = cl.Call(context.TODO(), "err", &test.NidOptNative{}, &resp)
	assert.Equal(t, testErr, err)
}

func TestBreacker(t *testing.T) {
	var (
		testData = test.NidOptNative{
			Field1:  -123,
			Field14: "dsfasdfasdf",
		}
		testErr = errors.New("Test Error!")
	)
	period := time.Millisecond

	srv := NewServer()
	srv.Handle("/ok", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		data, err := json.Marshal(&testData)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))
	srv.Handle("/err", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		Error(w, testErr, http.StatusInternalServerError)
	}))

	cl, err := NewClient("/", &JSONCodec{}, breaker.New(3, 1, period))
	assert.NoError(t, err)
	cl.SetRequestProcessor(func(req *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, req)
		return w.Result(), nil
	})

	var resp test.NidOptNative

	err = cl.Call(context.TODO(), "ok", &test.NidOptNative{}, &resp)
	assert.NoError(t, err)

	for i := 0; i < 3; i++ {
		err = cl.Call(context.TODO(), "err", &test.NidOptNative{}, &resp)
		assert.Equal(t, testErr, err)
	}

	err = cl.Call(context.TODO(), "err", &test.NidOptNative{}, &resp)
	assert.Equal(t, breaker.ErrBreakerOpen, err)
}

package testservice

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/qiwitech/graceful"
	"github.com/stretchr/testify/assert"
)

func TestSetter(t *testing.T) {
	mux := NewSrvSetterHandler(&underlyingSetter{}, &graceful.ProtobufCodec{})

	hcl, err := graceful.NewClient("/", &graceful.ProtobufCodec{}, nil)
	if err != nil {
		panic(err)
	}

	cl := NewSrvSetterHTTPClient(hcl)

	hcl.SetRequestProcessor(func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w.Result(), nil
	})

	resp, err := cl.Set(context.Background(), &ReqSet{Str: "str", Num: 10})
	if err != nil {
		t.Fatalf("call: %v", err)
	}

	if resp.Mod7 != 3 || resp.Upper != "STR" {
		t.Errorf("response: %v, want %v %v", resp, "STR", 3)
	}
}

func TestTimeoutOk(t *testing.T) {
	mux := NewSrvWaiterHandler(&underlyingWaiter{}, &graceful.ProtobufCodec{})

	hcl, err := graceful.NewClient("/", &graceful.ProtobufCodec{}, nil)
	if err != nil {
		panic(err)
	}

	cl := NewSrvWaiterHTTPClient(hcl)

	hcl.SetRequestProcessor(func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w.Result(), nil
	})

	_, err = cl.Wait(context.Background(), &ReqWait{Duration: time.Nanosecond.Nanoseconds()})
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	// ErrTimeout does not happed - so test successful
}

func TestTimeoutErr(t *testing.T) {
	mux := NewSrvWaiterHandler(&underlyingWaiter{}, &graceful.ProtobufCodec{})

	hcl, err := graceful.NewClient("/", &graceful.ProtobufCodec{}, nil)
	if err != nil {
		panic(err)
	}

	cl := NewSrvWaiterHTTPClient(hcl)

	hcl.SetRequestProcessor(func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w.Result(), nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	_, err = cl.Wait(ctx, &ReqWait{Duration: 100 * time.Millisecond.Nanoseconds()})
	assert.Error(t, err)
	assert.EqualError(t, err, context.DeadlineExceeded.Error())
}

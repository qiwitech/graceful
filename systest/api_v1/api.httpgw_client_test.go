// api.httpgw_client_test.go
package api_v1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/golang/mock/gomock"
	"github.com/qiwitech/graceful"
	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker(t *testing.T) {
	var (
		timeout      = time.Millisecond * 10
		errThreshold = 3
		okThreshold  = 1
		i            int
		called       bool
		testData     = []struct {
			fail      bool
			delay     time.Duration
			willBreak bool
		}{
			// request success: 4 times (tryLimit+1)
			{fail: false, delay: timeout / 5, willBreak: false},
			{fail: false, delay: timeout / 5, willBreak: false},
			{fail: false, delay: timeout / 5, willBreak: false},
			{fail: false, delay: timeout / 5, willBreak: false},
			// request fail: 2 times (tryLimit-1)
			{fail: true, delay: timeout / 5, willBreak: false},
			{fail: true, delay: timeout / 5, willBreak: false},
			// request success to close breaker: 1 times
			{fail: false, delay: timeout / 1, willBreak: false},
			// request fail: 3 times (tryLimit)
			{fail: true, delay: timeout / 5, willBreak: false},
			{fail: true, delay: timeout / 5, willBreak: false},
			{fail: true, delay: timeout / 5, willBreak: false},
			// other request after breaker open: X times
			{fail: false, delay: timeout * 2, willBreak: false},
			{fail: true, delay: timeout / 5, willBreak: false},
			{fail: true, delay: timeout / 5, willBreak: false},
			{fail: true, delay: timeout / 5, willBreak: false},
			{fail: false, delay: timeout / 5, willBreak: true},
		}
		expectCalls = 0
		failError   = errors.New("testError")
	)
	for _, tc := range testData {
		if !tc.fail && !tc.willBreak {
			expectCalls++
		}
	}
	// github.com/golang/mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// server mock implementation
	api := NewMockAPIInterface(ctrl)
	api.EXPECT().
		GetAccounts(gomock.Any(), gomock.Any()).
		Return(&AccountsResponse{}, nil).
		AnyTimes()

	// client
	client, srv, err := connectedAPIClient(api, breaker.New(errThreshold, okThreshold, timeout))
	if err != nil {
		t.Error(err)
	}
	client.SetRequestProcessor(func(req *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		called = true
		if testData[i].fail {
			graceful.Error(w, failError, http.StatusInternalServerError)
		} else {
			srv.ServeHTTP(w, req)
		}
		return w.Result(), nil
	})

	for i = 0; i < len(testData); i++ {
		called = false
		tc := testData[i]
		time.Sleep(tc.delay)
		_, err = client.GetAccounts(context.Background(), &AccountsRequest{})
		if tc.fail {
			assert.Equal(t, failError, err, "got in test case %d (%+v)", i, tc)
		}
		assert.Equal(t, tc.willBreak, !called, "got in test case %d (%+v) (error %v)", i, tc, err)
	}
}

func TestStatusCode(t *testing.T) {
	// github.com/golang/mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// server mock implementation
	api := NewMockAPIInterface(ctrl)
	// client
	client, _, err := connectedAPIClient(api, nil)
	if err != nil {
		t.Error(err)
	}

	const errorText = "ERROR"

	// slowly function must be called ClientTryLimit times only
	api.EXPECT().Transfer(gomock.Any(), gomock.Any()).
		Times(1).
		Return(&TransferResponse{}, errors.New(errorText))
	// call slowly function a lot of times
	_, err = client.Transfer(context.Background(), &TransferRequest{})
	if err == nil {
		t.Error("Client lose error")
	} else if err.Error() != errorText {
		t.Error("Client lose error text")
	}
}

func TestDeadline(t *testing.T) {
	// github.com/golang/mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// server mock implementation
	api := NewMockAPIInterface(ctrl)
	// client
	client, _, err := connectedAPIClient(api, nil)
	if err != nil {
		t.Error(err)
	}

	funcDelay := time.Millisecond * 100
	ctxTimeouts := []time.Duration{
		funcDelay * 2,
		funcDelay / 2,
	}

	// fast function must be called ClientTryLimit times only
	api.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, in *AccountsRequest) {
			time.Sleep(funcDelay)
		}).
		Return(&AccountsResponse{}, nil).
		AnyTimes() // Times(len(ctxTimeouts))

	for _, timeout := range ctxTimeouts {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		_, err := client.GetAccounts(ctx, &AccountsRequest{})
		// context deadline NOT exceeded
		if funcDelay < timeout && err != nil {
			t.Error(err)
		}
		// context deadline exceeded
		if funcDelay > timeout && err == nil {
			t.Error(err)
		}
		cancel()
	}

}

func connectedAPIClient(service APIInterface, b *breaker.Breaker) (*APIHTTPClient, *graceful.Server, error) {
	// fake listener
	l := graceful.NewFakeListener()
	http.DefaultTransport = l.DefaultTransport
	// server
	srv := graceful.NewServer()
	srv.Mount("/", NewAPIHandler(service, &graceful.JSONCodec{}))
	go srv.Serve(l)
	// client
	cl, err := graceful.NewClient("http://localhost", &graceful.JSONCodec{}, b)
	if err != nil {
		return nil, nil, err
	}
	c := NewAPIHTTPClient(cl)
	// ready
	return &c, srv, nil
}

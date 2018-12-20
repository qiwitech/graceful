// server_test.go
package graceful

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogo/protobuf/test"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {

	handler := NewHandler(&ProtobufCodec{},
		func() interface{} { return &test.NidOptNative{} },
		func(ctx context.Context, arg interface{}) (interface{}, error) {
			return arg, nil
		})

	req := httptest.NewRequest("POST", "/", nil)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusOK)

}

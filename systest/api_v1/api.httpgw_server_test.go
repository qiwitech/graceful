// api.httpgw_server_test.go
package api_v1

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qiwitech/graceful"
)

func TestHandlersCallService(t *testing.T) {
	// github.com/golang/mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// server mock implementation
	api := NewMockAPIInterface(ctrl)
	s := graceful.NewServer()
	s.Mount("/", NewAPIHandler(api, &graceful.JSONCodec{}))

	// data for tests
	testData := []struct {
		Path   string
		In     interface{}
		Expect func(interface{}, interface{}) *gomock.Call
	}{
		{"/GetAccounts", &AccountsRequest{}, api.EXPECT().GetAccounts},
		{"/GetAccountSettings", &AccountSettingsRequest{}, api.EXPECT().GetAccountSettings},
		{"/GetHistory", &HistoryRequest{}, api.EXPECT().GetHistory},
		{"/GetPrevHash", &PrevHashRequest{}, api.EXPECT().GetPrevHash},
		{"/GetStats", &StatsRequest{}, api.EXPECT().GetStats},
		{"/Transfer", &TransferRequest{}, api.EXPECT().Transfer},
		{"/UpdateSettings", &UpdateSettingsRequest{}, api.EXPECT().UpdateSettings},
	}
	// set expecting and send request
	for _, test := range testData {
		// set expecting
		test.Expect(gomock.Any(), test.In).Times(1)
		// send request
		body, err := json.Marshal(test.In)
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest("POST", test.Path, bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		s.ServeHTTP(resp, req)
	}
}

// pluto.httpgw_server_test.go
package pluto

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/mock/gomock"
	"github.com/qiwitech/graceful"
)

func TestHandlersCallProcessing(t *testing.T) {
	// github.com/golang/mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// server mock implementation
	proc := NewMockProcessingInterface(ctrl)
	s := graceful.NewServer()
	s.Mount("/", NewProcessingHandler(proc, &graceful.ProtobufCodec{}))

	// data for tests
	testData := []struct {
		Path   string
		In     interface{}
		Expect func(interface{}, interface{}) *gomock.Call
	}{
		{"/Transfer", &TransferRequest{}, proc.EXPECT().Transfer},
		{"/UpdateSettings", &UpdateSettingsRequest{}, proc.EXPECT().UpdateSettings},
	}
	// set expecting and send request
	for _, test := range testData {
		// set expecting
		test.Expect(gomock.Any(), test.In).Times(1)
		// send request
		body, err := proto.Marshal(test.In.(proto.Message))
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest("POST", test.Path, bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		s.ServeHTTP(resp, req)
	}
}

func TestHandlersCallStorage(t *testing.T) {
	// github.com/golang/mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// server mock implementation
	stor := NewMockStorageInterface(ctrl)
	s := graceful.NewServer()
	s.Mount("/", NewStorageHandler(stor, &graceful.ProtobufCodec{}))

	// data for tests
	testData := []struct {
		Path   string
		In     proto.Message
		Expect func(interface{}, interface{}) *gomock.Call
	}{
		{"/GetAccounts", &AccountsRequest{}, stor.EXPECT().GetAccounts},
		{"/GetAccountSettings", &AccountSettingsRequest{}, stor.EXPECT().GetAccountSettings},
		{"/GetHistory", &HistoryRequest{}, stor.EXPECT().GetHistory},
		{"/GetPrevHash", &PrevHashRequest{}, stor.EXPECT().GetPrevHash},
		{"/GetStats", &StatsRequest{}, stor.EXPECT().GetStats},
	}
	// set expecting and send request
	for _, test := range testData {
		// set expecting
		test.Expect(gomock.Any(), test.In).Times(1)
		// send request
		body, err := proto.Marshal(test.In)
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest("POST", test.Path, bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		s.ServeHTTP(resp, req)
	}
}

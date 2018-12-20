// fields_test.go
package api_v1

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qiwitech/graceful"
	"github.com/stretchr/testify/assert"
)

func TestAccountMarshaling(t *testing.T) {

	codecs := []graceful.Codec{
		&graceful.ProtobufCodec{},
		&graceful.JSONCodec{},
	}
	for _, c := range codecs {

		v0 := Account(0)
		buf, err := c.Marshal(&v0)
		assert.NoError(t, err, v0)

		v1 := v0
		err = c.Unmarshal(buf, &v1)
		assert.NoError(t, err, v1)
	}
}

func TestHashMarshaling(t *testing.T) {

	codecs := []graceful.Codec{
		&graceful.ProtobufCodec{},
		&graceful.JSONCodec{},
	}
	for _, c := range codecs {

		v0 := Hash("012345678901234567890123")
		buf, err := c.Marshal(&v0)
		assert.NoError(t, err)

		v1 := v0
		err = c.Unmarshal(buf, &v1)
		assert.NoError(t, err)

		assert.Equal(t, v0, v1)
	}
}

func TestUnmarshalJSONErrorType(t *testing.T) {
	var tr TransferRequest
	// тестовые пары с невалидными значениями
	testData := map[reflect.Type][]byte{
		reflect.TypeOf(tr.Sender):   []byte(`{"sender":"fffffffffffffffffffff"}`),
		reflect.TypeOf(tr.PrevHash): []byte(`{"prev_hash":"aaaa"}`),
	}
	// проверка типа возвращаемой ошибки (чтоб ошибки нормально читались клиентами)
	for fieldType, invalidRaw := range testData {
		err := json.Unmarshal(invalidRaw, &tr)
		if err == nil {
			t.Errorf("Bad testData for type %s", fieldType)
		} else if _, ok := err.(*json.UnmarshalTypeError); !ok {
			t.Errorf("%s.UnmarshalJSON() return invalid error type: %s", fieldType, reflect.TypeOf(err))
		}
	}
}

func TestTransferRequestJsonMarshaling(t *testing.T) {
	// эталонные пары

	testData := map[*TransferRequest][]byte{
		&TransferRequest{
			Sender:   0xFFFFFFFFFFFFFFFF,
			PrevHash: "8no7ah7YAj8r5AfJd/LQpw==",
		}: []byte(`{"sender":"ffffffffffffffff","prev_hash":"8no7ah7YAj8r5AfJd/LQpw=="}`),
		&TransferRequest{
			Sender:   0x1,
			PrevHash: "8no7ah7YAj8r5AfJd/LQpw==",
			Batch: []*Batch{
				&Batch{},
			},
		}: []byte(`{"sender":"1","prev_hash":"8no7ah7YAj8r5AfJd/LQpw==","batch":[{"receiver":"0"}]}`),
	}
	// сравнение
	for tr0, etalon := range testData {
		// encoding
		raw, err := json.Marshal(tr0)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(raw, etalon) {
			t.Logf("%+v >> %s", tr0, raw)
			t.Errorf("%s != %s", raw, etalon)
		}
		// decoding
		tr1 := &TransferRequest{}
		err = json.Unmarshal(raw, tr1)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(*tr0, *tr1) {
			t.Logf("%s >> %+v", raw, tr1)
			t.Errorf("%+v != %+v", *tr0, *tr1)
		}
	}
}

func TestAccountJsonMarshaling(t *testing.T) {
	// эталонные пары
	testData := map[Account][]byte{
		0:                  []byte(`"0"`),
		9:                  []byte(`"9"`),
		16:                 []byte(`"10"`),
		0xFFFFFFFFFFFFFFFF: []byte(`"ffffffffffffffff"`),
	}
	// сравнение
	for account, golden := range testData {

		raw, err := account.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(raw, golden) {
			t.Errorf("%s != %s", raw, golden)
		}
	}
}

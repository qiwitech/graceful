package graceful

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/gogo/protobuf/proto"
)

type Codec interface {
	Size(interface{}) int
	Marshal(interface{}) ([]byte, error)
	MarshalTo([]byte, interface{}) (int, error)
	Unmarshal([]byte, interface{}) error
	MIME() string
}

type ProtobufCodec struct {
}

func (c *ProtobufCodec) Size(v interface{}) int {
	if v == nil {
		return 0
	}
	return proto.Size(v.(proto.Message))
}

func (c *ProtobufCodec) Marshal(v interface{}) ([]byte, error) {
	if IsNilInterface(v) {
		return nil, nil
	}
	m, ok := v.(proto.Message)
	if !ok {
		panic("not supported")
	}

	data, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *ProtobufCodec) MarshalTo(buf []byte, v interface{}) (int, error) {
	if IsNilInterface(v) {
		return 0, nil
	}
	m, ok := v.(proto.Message)
	if !ok {
		panic("not supported")
	}

	data, err := proto.Marshal(m)
	if err != nil {
		return 0, nil
	}

	if len(data) > len(buf) {
		return 0, io.ErrShortBuffer
	}

	return copy(buf, data), nil
}

func (c *ProtobufCodec) Unmarshal(buf []byte, v interface{}) error {
	m, ok := v.(proto.Message)
	if !ok {
		panic("not supported")
	}

	return proto.Unmarshal(buf, m)
}

func (*ProtobufCodec) MIME() string {
	return "application/x-protobuf"
}

type JSONCodec struct {
}

func (c *JSONCodec) Size(v interface{}) int {
	if v == nil {
		return 0
	}
	// TODO(nik): optimize this
	data, err := json.Marshal(v)
	if err != nil {
		panic("marshal error")
	}
	return len(data)
}

func (c *JSONCodec) Marshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *JSONCodec) MarshalTo(buf []byte, v interface{}) (int, error) {
	if v == nil {
		return 0, nil
	}
	// TODO(nik): optimize this
	data, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	if len(data) > len(buf) {
		return 0, io.ErrShortBuffer
	}
	return copy(buf, data), nil
}

func (c *JSONCodec) Unmarshal(buf []byte, v interface{}) error {
	if len(buf) == 0 {
		return nil
	}
	return json.Unmarshal(buf, v)
}

func (*JSONCodec) MIME() string {
	return "application/json"
}

type CodecPack struct {
	Default   Codec
	Supported map[string]Codec
}

func (c *CodecPack) Codec(req http.Header) Codec {
	// TODO(nik): select codec by Header
	return c.Default
}

func IsNilInterface(v interface{}) bool {
	if v == nil {
		return true
	}
	r := reflect.ValueOf(v)
	return r.IsNil()
}

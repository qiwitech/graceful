package graceful

import (
	"testing"

	"github.com/gogo/protobuf/test"
	"github.com/stretchr/testify/assert"
)

var testObjs = []*test.NidOptNative{
	{},
	{
		Field1:  -123,
		Field14: "dsfasdfasdf",
	},
	nil,
}

func TestProtoMarshalTo(t *testing.T) {
	c := &ProtobufCodec{}

	for _, req := range testObjs {
		size := c.Size(req)
		buf := make([]byte, size)

		_, err := c.MarshalTo(buf, req)
		assert.NoError(t, err)

		var req2 test.NidOptNative
		err = c.Unmarshal(buf, &req2)
		assert.NoError(t, err)

		if req == nil {
			req = &test.NidOptNative{}
		}

		req.XXX_sizecache = 0
		req.XXX_unrecognized = nil
		req2.XXX_sizecache = 0
		req2.XXX_unrecognized = nil

		assert.Equal(t, req, &req2, "data: %x, test: %v", buf, req)

		if t.Failed() {
			break
		}
	}
}

func TestProtoMarshal(t *testing.T) {
	c := &ProtobufCodec{}

	for _, req := range testObjs {
		buf, err := c.Marshal(req)
		assert.NoError(t, err, "on obj %T %v", req, req)

		var req2 test.NidOptNative
		err = c.Unmarshal(buf, &req2)
		assert.NoError(t, err)

		if req == nil {
			req = &test.NidOptNative{}
		}

		req.XXX_sizecache = 0
		req.XXX_unrecognized = nil
		req2.XXX_sizecache = 0
		req2.XXX_unrecognized = nil

		assert.Equal(t, req, &req2, "data: %x", buf)

		if t.Failed() {
			break
		}
	}
}

func TestJSONMarshalTo(t *testing.T) {
	c := &JSONCodec{}

	for _, req := range testObjs {
		size := c.Size(req)
		buf := make([]byte, size)

		_, err := c.MarshalTo(buf, req)
		assert.NoError(t, err, "on test %T %+v", req, req)

		var req2 test.NidOptNative
		err = c.Unmarshal(buf, &req2)
		assert.NoError(t, err, "on test %T %+v", req, req)

		if req == nil {
			req = &test.NidOptNative{}
		}

		assert.Equal(t, req, &req2, "data: %x", buf)

		if t.Failed() {
			break
		}
	}
}

func TestJSONMarshal(t *testing.T) {
	c := &JSONCodec{}

	for _, req := range testObjs {
		buf, err := c.Marshal(req)
		assert.NoError(t, err)

		var req2 test.NidOptNative
		err = c.Unmarshal(buf, &req2)
		assert.NoError(t, err)

		if req == nil {
			req = &test.NidOptNative{}
		}

		assert.Equal(t, req, &req2, "data: %x", buf)

		if t.Failed() {
			break
		}
	}
}

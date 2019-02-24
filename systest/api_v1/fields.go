// fields.go
package api_v1

import (
	"encoding/binary"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/gogo/protobuf/proto"
)

// Номер аккаунта

type Account uint64

func (m *Account) Reset()                    { *m = Account(0) }
func (m *Account) String() string            { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()               {}
func (*Account) Descriptor() ([]byte, []int) { return fileDescriptorApi, []int{20} }

// 	proto.Marshaler
func (a *Account) Marshal() ([]byte, error) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(*a))
	return b, nil
}

// proto.Unmarshaler
func (a *Account) Unmarshal(b []byte) error {
	i := binary.LittleEndian.Uint64(b)
	a = (*Account)(&i)
	return nil
}

// json.Marshaler
func (a *Account) MarshalJSON() ([]byte, error) {
	str := `"` + strconv.FormatUint(uint64(*a), 16) + `"`
	return []byte(str), nil
}

// json.Unmarshaler
func (a *Account) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	val, err := strconv.ParseUint(str, 16, 64)
	if err != nil {
		jsonErr := &json.UnmarshalTypeError{
			Value: `string "` + str + `"`,
			Type:  reflect.TypeOf(*a),
		}
		return jsonErr
	}
	*a = Account(val)
	return nil
}

// Validatable
func (a *Account) Validate() error {
	return nil
}

func (a *Account) Size() int {
	return 8
}

// Hash

type Hash string

func (m *Hash) Reset()                    { *m = Hash("") }
func (m *Hash) String() string            { return proto.CompactTextString(m) }
func (*Hash) ProtoMessage()               {}
func (*Hash) Descriptor() ([]byte, []int) { return fileDescriptorApi, []int{20} }

// 	proto.Marshaler
func (h *Hash) Marshal() ([]byte, error) {
	return []byte(*h), nil
}

// proto.Unmarshaler
func (h *Hash) Unmarshal(b []byte) error {
	*h = Hash(b) //.(Hash)
	return nil
}

// json.Unmarshaler
func (h *Hash) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	if len(str) != 0 && len(str) != 24 {
		jsonErr := &json.UnmarshalTypeError{
			Value: `string "` + str + `"`,
			Type:  reflect.TypeOf(*h),
		}
		return jsonErr
	}
	*h = Hash(str)
	return nil
}

func (a *Hash) Validate() error {
	return nil
}

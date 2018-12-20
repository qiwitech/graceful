package testservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImplements(t *testing.T) {

	a := &SrvWaiterHTTPClient{}
	assert.Implements(t, (*SrvWaiterInterface)(nil), a)

	b := &SrvSetterHTTPClient{}
	assert.Implements(t, (*SrvSetterInterface)(nil), b)
}

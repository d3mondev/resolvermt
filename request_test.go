package fastdns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	request := newRequestDNS("www.google.com", TypeA)
	_, _, err := request.Send("8.8.8.8:53")

	assert.Nil(t, err)
}

package multidns

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := newRequestDNS("www.google.com", TypeA)
	_, _, err := request.Send("8.8.8.8:53")

	assert.Nil(t, err)
}

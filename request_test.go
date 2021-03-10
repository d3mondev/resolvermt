package multidns

import (
	"testing"

	mock_ratelimit "github.com/d3mondev/multidns/mocks/ratelimit"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLimiter := mock_ratelimit.NewMockLimiter(ctrl)
	mockLimiter.EXPECT().Take()

	resolver := newResolver("8.8.8.8:53", 10)
	resolver.(*resolverRateLimited).limiter = mockLimiter

	request := newRequestDNS("www.google.com", TypeA)
	_, _, err := request.Send(resolver.Get())

	assert.Nil(t, err)
}

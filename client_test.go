package multidns

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	resolvers := []string{"8.8.8.8"}

	client := New(resolvers, 3, 10, 5)

	assert.NotNil(t, client)
}

func TestRealSleeper(t *testing.T) {
	sleeper := defaultSleeper{}
	sleepFor := 1 * time.Millisecond

	start := time.Now()
	sleeper.Sleep(sleepFor)
	dur := time.Since(start)

	assert.True(t, dur >= sleepFor)
}

func TestRealNewSender(t *testing.T) {
	sender := realNewSender("test", TypeA)

	assert.NotNil(t, sender)
}

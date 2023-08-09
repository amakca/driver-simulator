package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsChanClosable(t *testing.T) {
	var ch chan struct{}

	// return false
	assert.False(t, IsChanClosable(ch))

	// default -> return true
	ch = make(chan struct{})
	assert.True(t, IsChanClosable(ch))

	// case -> return ok(true)
	go func() {
		ch <- struct{}{}
	}()
	time.Sleep(time.Millisecond)
	assert.True(t, IsChanClosable(ch))

	// default -> return true
	go func() {
		<-ch
	}()
	time.Sleep(time.Millisecond)
	assert.True(t, IsChanClosable(ch))
	ch <- struct{}{}

	// case -> return ok(false)
	close(ch)
	assert.False(t, IsChanClosable(ch))
}

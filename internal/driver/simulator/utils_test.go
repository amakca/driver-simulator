package driver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsChanClosable(t *testing.T) {
	t.Run("Nil channel should return false", func(t *testing.T) {
		var ch chan struct{}
		result := IsChanClosable(ch)
		assert.False(t, result)
	})

	t.Run("Closed channel should return false", func(t *testing.T) {
		ch := make(chan struct{})
		close(ch)
		result := IsChanClosable(ch)
		assert.False(t, result)
	})

	t.Run("Open channel should return true", func(t *testing.T) {
		ch := make(chan struct{})
		result := IsChanClosable(ch)
		assert.True(t, result)
	})
}

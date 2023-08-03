package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsChanClosable(t *testing.T) {
	t.Run("Nil channel", func(t *testing.T) {
		var ch chan struct{}
		result := IsChanClosable(ch)
		assert.False(t, result)
	})

	t.Run("Closed channel", func(t *testing.T) {
		ch := make(chan struct{})
		close(ch)
		result := IsChanClosable(ch)
		assert.False(t, result)
	})

	t.Run("Open channel", func(t *testing.T) {
		ch := make(chan struct{})
		go func() {
			ch <- struct{}{}
		}()
		time.Sleep(time.Millisecond)
		result := IsChanClosable(ch)
		assert.True(t, result)
	})

	t.Run("Busy channel", func(t *testing.T) {
		ch := make(chan struct{})
		go func() {
			<-ch
		}()
		time.Sleep(time.Millisecond)
		result := IsChanClosable(ch)
		assert.True(t, result)
	})
}

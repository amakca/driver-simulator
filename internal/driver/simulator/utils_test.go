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

/*
func TestSimulator_ChannelSwitch(t *testing.T) {
	t.Run("Creating new openCh1 channel if not closable", func(t *testing.T) {
		openCh1 := make(chan struct{})
		closeCh1 := make(chan struct{})
		channelSwitch(nil, openCh1, closeCh1)
		assert.NotNil(t, openCh1)
	})

	t.Run("Creating new openCh2 channel if not closable", func(t *testing.T) {
		openCh2 := make(chan struct{})
		closeCh1 := make(chan struct{})
		channelSwitch(openCh2, nil, closeCh1)
		assert.NotNil(t, openCh2)
	})

	t.Run("Closing closeCh1 channel if closable", func(t *testing.T) {
		openCh1 := make(chan struct{})
		closeCh1 := make(chan struct{})
		close(closeCh1)
		channelSwitch(openCh1, nil, closeCh1)
		_, ok := <-closeCh1
		assert.False(t, ok)
	})
}
*/

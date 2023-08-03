package generator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	g, err := NewRandGen("25ms:1.0:2.0")
	assert.NoError(t, err)

	err = g.Start()
	assert.NoError(t, err)
	err = g.Start()
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	value := g.Value()
	assert.True(t, value >= 1.0)

	err = g.Stop()
	assert.NoError(t, err)
	time.Sleep(50 * time.Millisecond)
	value = g.Value()
	time.Sleep(50 * time.Millisecond)
	valueAfter := g.Value()
	assert.NotEqual(t, value, valueAfter)

	err = g.Stop()
	assert.NoError(t, err)
	time.Sleep(50 * time.Millisecond)
	value = g.Value()
	time.Sleep(50 * time.Millisecond)
	valueAfter = g.Value()
	assert.Equal(t, value, valueAfter)

	err = g.Stop()
	assert.ErrorIs(t, err, ErrGenAlreadyStop)
}

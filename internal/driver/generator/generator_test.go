package generator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_OneSubs(t *testing.T) {
	g, err := NewRandGen("50ms:1.0:2.0")
	assert.NoError(t, err)

	err = g.Start()
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	value := g.ValueFloat32()

	assert.True(t, value >= 1.0)
	assert.True(t, value <= 2.0)

	err = g.Stop()
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	value = g.ValueFloat32()
	time.Sleep(100 * time.Millisecond)
	valueAfter := g.ValueFloat32()
	assert.Equal(t, value, valueAfter)

	err = g.Stop()
	assert.ErrorIs(t, err, errGenAlreadyStop)
}

func TestGenerator_TwoSubs(t *testing.T) {
	g, err := NewRandGen("30ms:1.0:2.0")
	assert.NoError(t, err)

	err = g.Start()
	assert.NoError(t, err)
	err = g.Start()
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	value := g.ValueFloat32()

	assert.True(t, value >= 1.0)
	assert.True(t, value <= 2.0)

	err = g.Stop()
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	value = g.ValueFloat32()
	time.Sleep(100 * time.Millisecond)
	valueAfter := g.ValueFloat32()
	assert.NotEqual(t, value, valueAfter)

	err = g.Stop()
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	value = g.ValueFloat32()
	time.Sleep(100 * time.Millisecond)
	valueAfter = g.ValueFloat32()
	assert.Equal(t, value, valueAfter)
}

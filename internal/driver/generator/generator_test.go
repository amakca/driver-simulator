package generator

import (
	u "practice/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	g, err := NewRandGen("25ms:2.0:1.0")
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), g.subs)
	assert.Equal(t, float32(0), g.value)
	assert.Equal(t, (time.Millisecond * 25), g.sampleRate)
	expSet := &random{
		high: 2.0,
		low:  1.0,
	}
	assert.Equal(t, expSet, g.valuer)

	t.Run("Start-stop", func(t *testing.T) {
		err = g.Start()
		assert.NoError(t, err)
		assert.Equal(t, uint32(1), g.subs)
		assert.True(t, u.IsChanClosable(g.done))

		time.Sleep(time.Millisecond * 30)
		assert.NotEqual(t, float32(1), g.value)

		err = g.Stop()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0), g.subs)
		assert.False(t, u.IsChanClosable(g.done))
	})

	t.Run("Subscription", func(t *testing.T) {
		assert.Equal(t, uint32(0), g.subs)

		g.Start()
		assert.Equal(t, uint32(1), g.subs)

		g.Start()
		assert.Equal(t, uint32(2), g.subs)

		g.Stop()
		assert.Equal(t, uint32(1), g.subs)

		g.Stop()
		assert.Equal(t, uint32(0), g.subs)

		err = g.Stop()
		assert.ErrorIs(t, err, ErrGenAlreadyStop)
	})

	t.Run("Value", func(t *testing.T) {
		g.value = 5
		assert.Equal(t, g.value, g.Value())
		assert.Equal(t, []byte{0x0, 0x0, 0xa0, 0x40}, g.ValueBytes())
	})
}

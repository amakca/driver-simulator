package driver

import (
	m "practice/internal/models"
	str "practice/internal/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimulator(t *testing.T) {
	gensettings := GeneralSettings{
		ProgramLiveTime: time.Hour * 5,
		GenOptimization: false,
	}

	set := m.DriverSettings{
		General: &gensettings,
		Tags:    map[m.DataID]m.Formatter{},
	}
	set.Tags[m.DataID(1)] = &TagSettings{
		PollTime:  25 * time.Millisecond,
		GenConfig: "rand:25ms:1.0:2.0",
	}
	set.Tags[m.DataID(2)] = &TagSettings{
		PollTime:  25 * time.Millisecond,
		GenConfig: "rand:25ms:1.0:2.0",
	}
	set.Tags[m.DataID(4)] = &TagSettings{
		PollTime:  35 * time.Millisecond,
		GenConfig: "rand:25ms:2.0:3.0",
	}

	settings3 := &TagSettings{
		PollTime:  25 * time.Millisecond,
		GenConfig: "rand:30ms:1.0:3.0",
	}

	str, _ := str.New()
	str.Create(1)
	str.Create(2)
	str.Create(3)
	str.Create(4)

	sim, err := New(set, str)
	assert.ErrorIs(t, err, m.ErrLiveTimeLong)
	gensettings.ProgramLiveTime = time.Minute * 5

	sim, err = New(set, str)
	assert.NoError(t, err)
	assert.NotNil(t, sim)

	t.Run("Correct working", func(t *testing.T) {
		err = sim.Run()
		assert.NoError(t, err)

		time.Sleep(50 * time.Millisecond)

		undo, err := sim.TagCreate(3, settings3)
		assert.NoError(t, err)
		time.Sleep(50 * time.Millisecond)
		err = undo()
		assert.NoError(t, err)

		err = sim.TagSetValue(2, []byte{1})
		assert.NoError(t, err)

		_, err = sim.TagDelete(4)
		assert.NoError(t, err)

		undo, err = sim.TagDelete(2)
		assert.NoError(t, err)
		err = undo()
		assert.NoError(t, err)

		err = sim.Stop()
		assert.NoError(t, err)
		time.Sleep(50 * time.Millisecond)
		err = sim.Run()
		assert.NoError(t, err)
		err = sim.Reset()
		assert.NoError(t, err)
		err = sim.Run()
		assert.NoError(t, err)
		err = sim.Reset()
		assert.NoError(t, err)
	})

	t.Run("Incorrect working", func(t *testing.T) {
		_, err = sim.TagCreate(2, set.Tags[2])
		assert.ErrorIs(t, err, m.ErrDataExists)

		_, err = sim.TagDelete(3)
		assert.ErrorIs(t, err, m.ErrDataNotFound)

		sim.Run()
		err = sim.Run()
		assert.ErrorIs(t, err, m.ErrAlreadyRunning)

		sim.Stop()
		err = sim.Stop()
		assert.ErrorIs(t, err, m.ErrAlreadyStopped)

		sim.Reset()
		err = sim.Stop()
		assert.ErrorIs(t, err, m.ErrNotWorking)

		sim.Close()
		err = sim.Close()
		assert.ErrorIs(t, err, m.ErrAlreadyClosed)

		err = sim.Run()
		assert.ErrorIs(t, err, m.ErrAlreadyClosed)
	})

}

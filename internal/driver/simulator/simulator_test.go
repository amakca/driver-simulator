package driver

import (
	m "practice/internal/models"
	"practice/internal/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimulator(t *testing.T) {
	generalSettings := GeneralSettings{
		ProgramLiveTime: time.Hour * 5,
		GenOptimization: true,
	}

	driverSettings := m.DriverSettings{
		General: &generalSettings,
		Tags:    map[m.DataID]m.Formatter{},
	}
	driverSettings.Tags[m.DataID(1)] = &TagSettings{
		PollTime:  30 * time.Millisecond,
		GenConfig: "rand:45ms:1.0:2.0",
	}
	driverSettings.Tags[m.DataID(2)] = &TagSettings{
		PollTime:  30 * time.Millisecond,
		GenConfig: "rand:35ms:1.0:3.0",
	}
	driverSettings.Tags[m.DataID(3)] = &TagSettings{
		PollTime:  35 * time.Millisecond,
		GenConfig: "rand:45ms:1.0:2.0",
	}

	storage, err := storage.New()
	assert.NoError(t, err)

	storage.Create(1)
	storage.Create(2)
	storage.Create(3)

	simulator, err := New(driverSettings, storage)
	assert.ErrorIs(t, err, m.ErrLiveTimeLong)

	generalSettings.ProgramLiveTime = time.Minute * 5
	simulator, err = New(driverSettings, storage)
	assert.NoError(t, err)

	t.Run("New", func(t *testing.T) {
		assert.NotNil(t, simulator.pollGroup)
		assert.NotNil(t, simulator.genManager)
		assert.NotNil(t, simulator.start)
		assert.Equal(t, m.READY, simulator.state)
		assert.Equal(t, generalSettings, simulator.generalSettings)
		assert.Equal(t, storage, simulator.storage)

		for k := range driverSettings.Tags {
			dr := driverSettings.Tags[k]
			sim := simulator.tagsSettings[k]
			assert.Equal(t, *dr.(*TagSettings), sim)
		}

		assert.Equal(t, simulator.pollGroup[30*time.Millisecond][1],
			simulator.pollGroup[35*time.Millisecond][3],
		)
		assert.NotEqual(t, simulator.pollGroup[30*time.Millisecond][1],
			simulator.pollGroup[30*time.Millisecond][2],
		)
	})

	t.Run("TagCreate", func(t *testing.T) {
		storage.Create(4)
		settings := &TagSettings{
			PollTime:  30 * time.Millisecond,
			GenConfig: "rand:30ms:1.0:3.0",
		}

		_, err = simulator.TagCreate(m.DataID(2), driverSettings.Tags[m.DataID(2)])
		assert.ErrorIs(t, err, m.ErrDataExists)

		undo, err := simulator.TagCreate(m.DataID(4), settings)
		assert.NoError(t, err)
		assert.Contains(t, simulator.tagsSettings, m.DataID(4))
		assert.Contains(t, simulator.pollGroup[settings.PollTime], m.DataID(4))

		assert.NoError(t, undo())
		assert.NotContains(t, simulator.tagsSettings, m.DataID(4))
		assert.NotContains(t, simulator.pollGroup[settings.PollTime], m.DataID(4))
	})

	t.Run("TagDelete", func(t *testing.T) {
		_, err = simulator.TagDelete(m.DataID(4))
		assert.ErrorIs(t, err, m.ErrDataNotFound)

		undo, err := simulator.TagDelete(m.DataID(2))
		assert.NoError(t, err)
		assert.NotContains(t, simulator.tagsSettings, m.DataID(2))
		assert.NotContains(t, simulator.pollGroup[30*time.Millisecond], m.DataID(2))

		assert.NoError(t, undo())
		assert.Contains(t, simulator.tagsSettings, m.DataID(2))
		assert.Contains(t, simulator.pollGroup[30*time.Millisecond], m.DataID(2))

	})

	t.Run("TagSetValue", func(t *testing.T) {
		assert.NoError(t, simulator.TagSetValue(2, []byte{0, 0, 0, 40}))
		assert.Equal(t, []byte{0, 0, 0, 40},
			simulator.pollGroup[30*time.Millisecond][2].ValueBytes(),
		)

		assert.ErrorIs(t, simulator.TagSetValue(5, []byte{0, 0, 0, 40}),
			m.ErrDataNotFound)
	})

	t.Run("Settings", func(t *testing.T) {
		listSettings := simulator.Settings()
		assert.Equal(t, driverSettings, listSettings)
	})

	t.Run("Service", func(t *testing.T) {
		assert.ErrorIs(t, simulator.Stop(), m.ErrNotWorking)

		assert.NoError(t, simulator.Run())
		assert.ErrorIs(t, simulator.Run(), m.ErrAlreadyRunning)
		assert.Equal(t, m.RUNNING, simulator.state)

		time.Sleep(time.Millisecond * 50)

		assert.NoError(t, simulator.Stop())
		assert.ErrorIs(t, simulator.Stop(), m.ErrAlreadyStopped)
		assert.Equal(t, m.STOPPED, simulator.state)

		assert.NoError(t, simulator.Run())
		assert.Equal(t, m.RUNNING, simulator.state)

		assert.NoError(t, simulator.Reset())
		assert.Equal(t, m.READY, simulator.state)

		assert.NoError(t, simulator.Run())
		assert.Equal(t, m.RUNNING, simulator.state)

		_, err := simulator.TagDelete(m.DataID(3))
		assert.NoError(t, err)
		time.Sleep(time.Millisecond * 50)

		assert.NoError(t, simulator.Close())
		assert.ErrorIs(t, simulator.Close(), m.ErrProgramClosed)
		assert.ErrorIs(t, simulator.Stop(), m.ErrProgramClosed)
		assert.ErrorIs(t, simulator.Run(), m.ErrProgramClosed)
		assert.Equal(t, m.CLOSED, simulator.state)
	})

}

package storage

import (
	"testing"

	m "practice/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestTagStorage(t *testing.T) {
	storage, err := New()
	assert.NoError(t, err)
	id := m.DataID(1)

	t.Run("Create", func(t *testing.T) {
		undo, err := storage.Create(id)
		assert.NoError(t, err)
		assert.NotNil(t, undo)
		assert.Equal(t, 1, len(storage.data))

		_, err = storage.Create(id)
		assert.ErrorIs(t, err, m.ErrDataExists)

		err = undo()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(storage.data))
	})

	t.Run("Read", func(t *testing.T) {
		datapoint, err := storage.Read(11)
		assert.Equal(t, m.Datapoint{}, datapoint)
		assert.ErrorIs(t, err, m.ErrDataNotFound)

		datapoint = m.Datapoint{
			Value:     []byte{1},
			Timestamp: 2,
			Quality:   m.GOOD,
		}
		storage.data[id] = &datapoint

		datapoint, err = storage.Read(id)
		assert.NoError(t, err)
		assert.Equal(t, storage.data[id], &datapoint)
	})

	t.Run("Update", func(t *testing.T) {
		datapoint := &m.Datapoint{
			Value:     []byte{1},
			Timestamp: 2,
			Quality:   m.GOOD,
		}

		undo, err := storage.Update(11, *datapoint)
		assert.ErrorIs(t, err, m.ErrDataNotFound)

		storage.data[id] = &m.Datapoint{}

		undo, err = storage.Update(id, *datapoint)
		assert.NoError(t, err)
		assert.NotNil(t, undo)
		assert.Equal(t, storage.data[id], datapoint)

		err = undo()
		assert.NoError(t, err)
		assert.Equal(t, storage.data[id], &m.Datapoint{})
	})

	t.Run("Delete", func(t *testing.T) {
		_, err = storage.Delete(11)
		assert.ErrorIs(t, err, m.ErrDataNotFound)

		storage.data[id] = &m.Datapoint{}

		undo, err := storage.Delete(id)
		assert.NoError(t, err)
		assert.NotNil(t, undo)
		assert.Equal(t, 0, len(storage.data))

		err = undo()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(storage.data))
	})

	t.Run("UpdateValue", func(t *testing.T) {
		_, err = storage.UpdateValue(11, []byte{2})
		assert.ErrorIs(t, err, m.ErrDataNotFound)

		storage.data[id] = &m.Datapoint{}

		undo, err := storage.UpdateValue(id, []byte{1})
		assert.NoError(t, err)
		assert.NotNil(t, undo)
		assert.Equal(t, []byte{1}, storage.data[id].Value)

		err = undo()
		assert.NoError(t, err)
		assert.Equal(t, []byte(nil), storage.data[id].Value)
	})

	t.Run("UpdateQuality", func(t *testing.T) {
		_, err = storage.UpdateQuality(11, m.BAD)
		assert.ErrorIs(t, err, m.ErrDataNotFound)

		storage.data[id] = &m.Datapoint{}

		undo, err := storage.UpdateQuality(id, m.BAD)
		assert.NoError(t, err)
		assert.NotNil(t, undo)
		assert.Equal(t, m.BAD, storage.data[id].Quality)

		err = undo()
		assert.NoError(t, err)
		assert.Equal(t, m.UNCERTAIN, storage.data[id].Quality)
	})

	t.Run("List", func(t *testing.T) {
		storage.data[2] = &m.Datapoint{
			Value:     []byte{2},
			Timestamp: 2,
			Quality:   m.GOOD,
		}
		storage.data[3] = &m.Datapoint{
			Value:     []byte{3},
			Timestamp: 3,
			Quality:   m.BAD,
		}
	})

	listTags := storage.List()
	for k, v := range storage.data {
		assert.Equal(t, listTags[k], *v)
	}

}

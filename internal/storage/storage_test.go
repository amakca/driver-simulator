package storage

import (
	"testing"

	m "practice/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestTagStorage(t *testing.T) {
	str, err := New()
	assert.NoError(t, err)

	// Create
	id := m.DataID(1)
	undo, err := str.Create(id)
	assert.NoError(t, err)
	assert.NotNil(t, undo)

	_, err = str.Create(id)
	assert.ErrorIs(t, err, errDataExist)

	err = undo()
	assert.NoError(t, err)

	_, err = str.Read(id)
	assert.ErrorIs(t, err, errDataNotFound)

	// Update
	datapoint := m.Datapoint{Quality: m.Good}
	undo, err = str.Update(id, datapoint)
	assert.ErrorIs(t, err, errDataNotFound)

	_, err = str.Create(id)
	assert.NoError(t, err)

	undo, err = str.Update(id, datapoint)
	assert.NoError(t, err)
	assert.NotNil(t, undo)

	err = undo()
	assert.NoError(t, err)

	read, err := str.Read(id)
	assert.NoError(t, err)
	assert.Equal(t, read.Quality, m.Uncertain)

	// Delete
	undo, err = str.Delete(id)
	assert.NoError(t, err)
	assert.NotNil(t, undo)

	_, err = str.Delete(id)
	assert.ErrorIs(t, err, errDataNotFound)

	err = undo()
	assert.NoError(t, err)

	_, err = str.Read(id)
	assert.NoError(t, err)

	// UpdateValue
	undo, err = str.UpdateValue(id, []byte{1})
	assert.NoError(t, err)
	assert.NotNil(t, undo)

	err = undo()
	assert.NoError(t, err)

	read, err = str.Read(id)
	assert.Nil(t, read.Value)

	_, err = str.UpdateValue(2, []byte{2})
	assert.ErrorIs(t, err, errDataNotFound)

	// UpdateQuality
	undo, err = str.UpdateQuality(id, m.Bad)
	assert.NoError(t, err)
	assert.NotNil(t, undo)
	read, err = str.Read(id)
	assert.Equal(t, read.Quality, m.Bad)

	err = undo()
	assert.NoError(t, err)

	read, err = str.Read(id)
	assert.Equal(t, read.Quality, m.Uncertain)

	_, err = str.UpdateQuality(2, m.Bad)
	assert.ErrorIs(t, err, errDataNotFound)

	// List
	str.Create(m.DataID(2))
	str.UpdateValue(m.DataID(2), []byte{2})

	data := str.List()
	assert.Equal(t, len(data), 2)
	assert.Contains(t, data, id)
	assert.Contains(t, data, m.DataID(2))

}

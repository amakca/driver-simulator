package driver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTagSettings_String(t *testing.T) {
	tagSet := &TagSettings{
		PollTime: 5 * time.Second,
		Settings: "example",
	}

	expected := "5sexample"
	actual := tagSet.String()

	assert.Equal(t, expected, actual)
}

func TestTagSettings_BytesJSON(t *testing.T) {
	tagSet := &TagSettings{
		PollTime: time.Microsecond,
		Settings: "example",
	}

	expected := []byte(`{"poll-time":1000,"settings":"example"}`)
	actual, err := tagSet.BytesJSON()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseTags_JSON(t *testing.T) {
	input := []byte(`{"poll-time":1000,"settings":"example"}`)

	expected := TagSettings{
		PollTime: time.Microsecond,
		Settings: "example",
	}
	actual, err := parseTags(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseTags_String(t *testing.T) {
	input := "5s:rand:1s:1.0:2.0"

	expected := TagSettings{
		PollTime: 5 * time.Second,
		Settings: "rand:1s:1.0:2.0",
	}
	actual, err := parseTags(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseTags_Struct(t *testing.T) {
	tagSet := &TagSettings{
		PollTime: 5 * time.Second,
		Settings: "rand:1s:1.0:2.0",
	}

	expected := *tagSet
	actual, err := parseTags(tagSet)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseTags_InvalidInput(t *testing.T) {
	input := "invalid_input"

	_, err := parseTags(input)

	assert.Error(t, err)
	assert.EqualError(t, err, errInvalidSettings.Error())
}

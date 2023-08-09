package driver

import (
	m "practice/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTagsSettings(t *testing.T) {
	tagSettings := &TagSettings{
		PollTime:  30 * time.Millisecond,
		GenConfig: "rand:1s:1.0:2.0",
	}

	t.Run("String", func(t *testing.T) {
		expected := "30ms:rand:1s:1.0:2.0"
		actual := tagSettings.String()
		assert.Equal(t, expected, actual)
	})

	t.Run("BytesJSON", func(t *testing.T) {
		expected := []byte(`{"poll-time":30000000,"generator-config":"rand:1s:1.0:2.0"}`)
		actual, err := tagSettings.BytesJSON()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("parseTags", func(t *testing.T) {
		input := "invalid_input"
		_, err := parseTags(input)
		assert.EqualError(t, err, m.ErrInvalidSettings.Error())

		input = "10ms:config"
		_, err = parseTags(input)
		assert.EqualError(t, err, m.ErrPollTimeSmall.Error())

		input = tagSettings.String()
		actual, err := parseTags(input)
		assert.NoError(t, err)
		assert.Equal(t, *tagSettings, actual)
	})

	t.Run("parseTagsJSON", func(t *testing.T) {
		input := []byte(`{"poll-time":30000000,"generator-config":"rand:1s:1.0:2.0"}`)
		actual, err := parseTags(input)
		assert.NoError(t, err)
		assert.Equal(t, *tagSettings, actual)
	})

	t.Run("parseTagsString", func(t *testing.T) {
		input := "30ms:rand:1s:1.0:2.0"
		actual, err := parseTags(input)
		assert.NoError(t, err)
		assert.Equal(t, *tagSettings, actual)
	})

	t.Run("parseTagsStruct", func(t *testing.T) {
		actual, err := parseTags(tagSettings)
		assert.NoError(t, err)
		assert.Equal(t, *tagSettings, actual)
	})
}

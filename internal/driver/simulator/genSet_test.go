package driver

import (
	m "practice/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGeneralSettings(t *testing.T) {
	generalSettings := &GeneralSettings{
		ProgramLiveTime: time.Microsecond,
		GenOptimization: true,
	}

	t.Run("String", func(t *testing.T) {
		expected := "1µs:true"
		actual := generalSettings.String()
		assert.Equal(t, expected, actual)
	})

	t.Run("BytesJSON", func(t *testing.T) {
		expected := []byte(`{"program-live-time":1000,"use-generator-optimization":true}`)
		actual, err := generalSettings.BytesJSON()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("parseGeneral", func(t *testing.T) {
		input := "invalid_input"
		_, err := parseGeneral(input)
		assert.EqualError(t, err, m.ErrInvalidSettings.Error())

		input = "2h:true"
		_, err = parseGeneral(input)
		assert.EqualError(t, err, m.ErrLiveTimeLong.Error())

		input = generalSettings.String()
		actual, err := parseGeneral(input)
		assert.NoError(t, err)
		assert.Equal(t, *generalSettings, actual)
	})

	t.Run("parseGeneralJSON", func(t *testing.T) {
		input := []byte(`{"program-live-time":1000,"use-generator-optimization":true}`)
		actual, err := parseGeneral(input)
		assert.NoError(t, err)
		assert.Equal(t, *generalSettings, actual)
	})

	t.Run("parseGeneralString", func(t *testing.T) {
		input := "1µs:true"
		actual, err := parseGeneral(input)
		assert.NoError(t, err)
		assert.Equal(t, *generalSettings, actual)
	})

	t.Run("parseGeneralStruct", func(t *testing.T) {
		actual, err := parseGeneral(generalSettings)
		assert.NoError(t, err)
		assert.Equal(t, *generalSettings, actual)
	})
}

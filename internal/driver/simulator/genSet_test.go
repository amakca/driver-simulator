package driver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGeneralSettings_String(t *testing.T) {
	genSet := &GeneralSettings{
		MaxLiveTime:   time.Minute,
		UseGenManager: true,
	}

	expected := "1m0s true"
	actual := genSet.String()

	assert.Equal(t, expected, actual)
}

func TestGeneralSettings_BytesJSON(t *testing.T) {
	genSet := &GeneralSettings{
		MaxLiveTime:   time.Microsecond,
		UseGenManager: true,
	}

	expected := []byte(`{"max-live-time":1000,"flag-generator-manager":true}`)
	actual, err := genSet.BytesJSON()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseGeneral_JSON(t *testing.T) {
	input := []byte(`{"max-live-time":1000,"flag-generator-manager":true}`)

	expected := GeneralSettings{
		MaxLiveTime:   time.Microsecond,
		UseGenManager: true,
	}
	actual, err := parseGeneral(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseGeneral_String(t *testing.T) {
	input := "1ms:true"

	expected := GeneralSettings{
		MaxLiveTime:   time.Millisecond,
		UseGenManager: true,
	}
	actual, err := parseGeneral(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseGeneral_Struct(t *testing.T) {
	genSet := &GeneralSettings{
		MaxLiveTime:   time.Minute,
		UseGenManager: true,
	}

	expected := *genSet
	actual, err := parseGeneral(genSet)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseGeneral_InvalidInput(t *testing.T) {
	input := "invalid_input"

	_, err := parseGeneral(input)

	assert.Error(t, err)
	assert.EqualError(t, err, errInvalidSettings.Error())
}

func TestParseGeneral_LongLiveTime(t *testing.T) {
	input := "2h:true"

	_, err := parseGeneral(input)

	assert.Error(t, err, "An error should occur for long live time")
	assert.EqualError(t, err, errLiveTimeLong.Error(), "Error should match")
}

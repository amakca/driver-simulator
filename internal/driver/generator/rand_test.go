package generator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseRandSettings(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       random
		wantSample time.Duration
		wantErr    bool
	}{
		{
			name:       "valid",
			input:      "1s:1.0:2.0",
			want:       random{1.0, 2.0},
			wantSample: 1 * time.Second,
			wantErr:    false,
		},
		{
			name:       "not enough parts",
			input:      "50ms:1.0",
			want:       random{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
		{
			name:       "sampleRate too small",
			input:      "1ms:1.0:2.0",
			want:       random{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
		{
			name:       "invalid delimiter",
			input:      "1s;1.0;2.0",
			want:       random{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotSample, err := parseRandom(tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantSample, gotSample)
			if tt.wantErr {
				assert.Error(t, err)

			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewRandGen(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantSettings random
		wantErr      bool
	}{
		{
			name:  "valid",
			input: "1s:1.0:2.0",
			wantSettings: random{
				low:  1.0,
				high: 2.0,
			},
			wantErr: false,
		},
		{
			name:         "invalid input",
			input:        "foo:2.0:1s",
			wantSettings: random{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRandGen(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValue(t *testing.T) {
	settings := random{
		low:  1.0,
		high: 2.0,
	}

	for i := 0; i < 100; i++ {
		value := settings.value()
		assert.True(t, value >= 1.0)
		assert.True(t, value <= 2.0)
	}
}

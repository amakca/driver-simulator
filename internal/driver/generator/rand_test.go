package generator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRandSettings(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       randSettings
		wantSample time.Duration
		wantErr    bool
	}{
		{
			name:       "valid",
			input:      "1s:1.0:2.0",
			want:       randSettings{1.0, 2.0},
			wantSample: 1 * time.Second,
			wantErr:    false,
		},
		{
			name:       "not enough parts",
			input:      "50ms:1.0",
			want:       randSettings{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
		{
			name:       "sampleRate too small",
			input:      "1ms:1.0:2.0",
			want:       randSettings{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
		{
			name:       "invalid delimiter",
			input:      "1s;1.0;2.0",
			want:       randSettings{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotSample, err := parseRandSettings(tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantSample, gotSample)
			if tt.wantErr {
				require.Error(t, err)

			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewRandGen(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantSettings randSettings
		wantErr      bool
	}{
		{
			name:  "valid",
			input: "1s:1.0:2.0",
			wantSettings: randSettings{
				low:  1.0,
				high: 2.0,
			},
			wantErr: false,
		},
		{
			name:         "invalid input",
			input:        "foo:2.0:1s",
			wantSettings: randSettings{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRandGen(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValue(t *testing.T) {
	settings := randSettings{
		low:  1.0,
		high: 2.0,
	}

	for i := 0; i < 100; i++ {
		value := settings.value()
		assert.True(t, value >= 1.0)
		assert.True(t, value <= 2.0)
	}
}

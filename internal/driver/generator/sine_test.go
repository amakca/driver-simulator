package generator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseSineSettings(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       sineSettings
		wantSample time.Duration
		wantErr    bool
	}{
		{
			name:       "valid",
			input:      "1s:1.0:2.0",
			want:       sineSettings{1.0, 2.0},
			wantSample: 1 * time.Second,
			wantErr:    false,
		},
		{
			name:       "not enough parts",
			input:      "50ms:1.0",
			want:       sineSettings{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
		{
			name:       "sampleRate too small",
			input:      "1ms:1.0:2.0",
			want:       sineSettings{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
		{
			name:       "invalid delimiter",
			input:      "1s;1.0;2.0",
			want:       sineSettings{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotSample, err := parseSineSettings(tt.input)
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

func TestNewSineGen(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantSettings sineSettings
		wantErr      bool
	}{
		{
			name:  "valid",
			input: "1s:1.0:2.0",
			wantSettings: sineSettings{
				amplitude: 1.0,
				frequency: 2.0,
			},
			wantErr: false,
		},
		{
			name:         "invalid input",
			input:        "foo:2.0:1s",
			wantSettings: sineSettings{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSineGen(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

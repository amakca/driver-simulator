package generator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSawSettings(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       sawSettings
		wantSample time.Duration
		wantErr    bool
	}{
		{
			name:       "valid",
			input:      "1s:1.0:2.0",
			want:       sawSettings{1.0, 2.0},
			wantSample: 1 * time.Second,
			wantErr:    false,
		},
		{
			name:       "not enough parts",
			input:      "50ms:1.0",
			want:       sawSettings{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
		{
			name:       "sampleRate too small",
			input:      "1ms:1.0:2.0",
			want:       sawSettings{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
		{
			name:       "invalid delimiter",
			input:      "1s;1.0;2.0",
			want:       sawSettings{},
			wantSample: 0 * time.Second,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotSample, err := parseSawSettings(tt.input)
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

func TestNewSawGen(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantSettings sawSettings
		wantErr      bool
	}{
		{
			name:  "valid",
			input: "1s:1.0:2.0",
			wantSettings: sawSettings{
				amplitude: 1.0,
				frequency: 2.0,
			},
			wantErr: false,
		},
		{
			name:         "invalid input",
			input:        "foo:2.0:1s",
			wantSettings: sawSettings{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSawGen(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

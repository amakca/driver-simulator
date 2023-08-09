package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManagerNew(t *testing.T) {
	genManager, err := CreateManager()
	assert.NotNil(t, genManager.list)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		input   string
		flag    bool
		len     int
		wantErr bool
	}{
		{
			name:    "Without flag, existing generator",
			input:   "rand:1s:1.0:2.0",
			flag:    false,
			len:     0,
			wantErr: false,
		},
		{
			name:    "With flag, existing generator",
			input:   "rand:1s:1.0:2.0",
			flag:    true,
			len:     1,
			wantErr: false,
		},
		{
			name:    "With flag, non-existing generator",
			input:   "rand:1s:1.0:3.0",
			flag:    true,
			len:     2,
			wantErr: false,
		},
		{
			name:    "With invalid config",
			input:   "invalidConfig",
			flag:    false,
			len:     2,
			wantErr: true,
		},
		{
			name:    "With flag, two existing generators",
			input:   "rand:1s:1.0:3.0",
			flag:    true,
			len:     2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := genManager.New(tt.input, tt.flag)
			assert.Equal(t, tt.len, len(genManager.list))
			if tt.wantErr {
				assert.Nil(t, gen)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, gen)
				assert.NoError(t, err)
				if tt.flag {
					assert.NotNil(t, genManager.list[tt.input])
				}
			}
		})
	}
}

func TestGenManagerParse(t *testing.T) {
	genManager, _ := CreateManager()

	t.Run("Valid config", func(t *testing.T) {
		genType, genCfg, err := genManager.parseConfig("rand:1s:3.0:1.0")
		assert.NoError(t, err)
		assert.Equal(t, "rand", genType)
		assert.Equal(t, "1s:3.0:1.0", genCfg)
	})

	t.Run("Invalid config", func(t *testing.T) {
		genType, genCfg, err := genManager.parseConfig("invalidConfig")
		assert.Error(t, err)
		assert.Equal(t, "", genType)
		assert.Equal(t, "", genCfg)
	})

	t.Run("Select sine generator", func(t *testing.T) {
		gen, err := genManager.selectGenType("sine", "1s:1.0:3.0")
		assert.Equal(t, &sine{
			amplitude: 1.0,
			frequency: 3.0,
		}, gen.valuer)
		assert.NoError(t, err)
	})

	t.Run("Select saw generator", func(t *testing.T) {
		gen, err := genManager.selectGenType("saw", "1s:1.0:3.0")
		assert.Equal(t, &saw{
			amplitude: 1.0,
			frequency: 3.0,
		}, gen.valuer)
		assert.NoError(t, err)
	})

	t.Run("Select rand generator", func(t *testing.T) {
		gen, err := genManager.selectGenType("rand", "1s:3.0:1.0")
		assert.Equal(t, &random{
			high: 3.0,
			low:  1.0,
		}, gen.valuer)
		assert.NoError(t, err)
	})

	t.Run("Select non-existing generator", func(t *testing.T) {
		gen, err := genManager.selectGenType("invalidGen", "1s:1.0:3.0")
		assert.Error(t, err)
		assert.Nil(t, gen)
	})
}

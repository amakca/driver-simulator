package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGenManager(t *testing.T) {
	genManager := CreateGenManager()

	assert.NotNil(t, genManager)
	assert.Empty(t, genManager.listGen)
}

func TestGenManager_New(t *testing.T) {
	genManager := CreateGenManager()

	// Test case 1: New generator with flag set to false
	gen, err := genManager.New("rand:1s:1.0:2.0", false)
	assert.NoError(t, err)
	assert.NotNil(t, gen)
	assert.Equal(t, 0, len(genManager.listGen))

	// Test case 2: New generator with flag set to true and existing generator
	gen, err = genManager.New("rand:1s:1.0:2.0", true)
	assert.NoError(t, err)
	assert.NotNil(t, gen)
	assert.Equal(t, 1, len(genManager.listGen))

	// Test case 3: New generator with flag set to true and non-existing generator
	gen, err = genManager.New("rand:1s:1.0:3.0", true)
	assert.NoError(t, err)
	assert.NotNil(t, gen)
	assert.Equal(t, 2, len(genManager.listGen))

	// Test case 4: New generator with invalid config
	gen, err = genManager.New("invalidConfig", false)
	assert.Error(t, err)
	assert.Nil(t, gen)
	assert.Equal(t, 2, len(genManager.listGen))

	// Test case 5: New generator with flag set to true and existing generator
	gen, err = genManager.New("rand:1s:1.0:3.0", true)
	assert.NoError(t, err)
	assert.NotNil(t, gen)
	assert.Equal(t, 2, len(genManager.listGen))
}

func TestGenManager_parseConfig(t *testing.T) {
	genManager := CreateGenManager()

	// Test case 1: Valid config
	genType, genCfg, err := genManager.parseConfig("rand:1s:1.0:3.0")
	assert.NoError(t, err)
	assert.Equal(t, "rand", genType)
	assert.Equal(t, "1s:1.0:3.0", genCfg)

	// Test case 2: Invalid config
	genType, genCfg, err = genManager.parseConfig("invalidConfig")
	assert.Error(t, err)
	assert.Equal(t, "", genType)
	assert.Equal(t, "", genCfg)
}

func TestGenManager_selectGenType(t *testing.T) {
	genManager := CreateGenManager()

	// Test case 1: Sine generator
	gen, err := genManager.selectGenType("sine", "1s:1.0:3.0")
	assert.NoError(t, err)
	assert.NotNil(t, gen)

	// Test case 2: Saw generator
	gen, err = genManager.selectGenType("saw", "1s:1.0:3.0")
	assert.NoError(t, err)
	assert.NotNil(t, gen)

	// Test case 3: Random generator
	gen, err = genManager.selectGenType("rand", "1s:1.0:3.0")
	assert.NoError(t, err)
	assert.NotNil(t, gen)

	// Test case 4: Invalid generator type
	gen, err = genManager.selectGenType("invalidGen", "1s:1.0:3.0")
	assert.Error(t, err)
	assert.Nil(t, gen)
}

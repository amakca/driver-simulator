package generator

import (
	"errors"
	"time"
)

var (
	ErrInvalidSettings   = errors.New("invalid settings format")
	ErrGenTypeNotFound   = errors.New("generator type not found")
	ErrSampleRateSmall   = errors.New("sample rate time is too small")
	ErrGenAlreadyStopped = errors.New("generator already stopped")
)

const (
	MAX_SAMPLE_RATE time.Duration = time.Millisecond * 25

	SINE_GENERATOR = "sine"
	SAW_GENERATOR  = "saw"
	RAND_GENERATOR = "rand"
)

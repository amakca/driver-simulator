package generator

import (
	"errors"
	"time"
)

var (
	errInvalidSettings = errors.New("invalid settings format")
	errGenTypeNotFound = errors.New("generator type not found")
	errPrescallerSmall = errors.New("prescaler is too small")
	errGenAlreadyStop  = errors.New("generator already stopped")
)

const (
	maxPrescaler time.Duration = time.Millisecond * 25
	delimiter    string        = ":"

	sineGen = "sine"
	sawGen  = "saw"
	randGen = "rand"
)

package driver

import (
	"errors"
	"time"
)

var (
	errDataExist       = errors.New("data with id already exist")
	errDataNotFound    = errors.New("data with id not found")
	errAlreadyRunning  = errors.New("program already running")
	errAlreadyClosed   = errors.New("program already closed")
	errAlreadyStopped  = errors.New("program already stopped")
	errUnknownState    = errors.New("program state unknown")
	errProgramNotReady = errors.New("program not ready")
	errNotWorking      = errors.New("program not working")
	errLiveTimeLong    = errors.New("live time is too long")
	errPrescallerSmall = errors.New("prescaler is too small")
	errInvalidSettings = errors.New("invalid settings format")
)

const (
	maxPrescaler = time.Millisecond * 25
	MaxLiveTime  = time.Hour

	delimiter  = ":"
	configFile = "config.json"
)

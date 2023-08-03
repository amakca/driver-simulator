package models

import (
	"time"

	"github.com/pkg/errors"
)

const (
	UNCERTAIN QualityState = iota
	GOOD
	BAD
)

const (
	READY uint8 = iota
	RUNNING
	STOPPED
	RESET
	CLOSED
)

const (
	MIN_POLL_TIME = time.Millisecond * 25
	MAX_LIVE_TIME = time.Hour

	DELIMITER   = ":"
	CONFIG_FILE = "config.json"
)

var (
	ErrDataExists   = errors.New("data with id already exists: ")
	ErrDataNotFound = errors.New("data with id not found :")

	ErrAlreadyRunning    = errors.New("program already running")
	ErrAlreadyClosed     = errors.New("program already closed")
	ErrAlreadyStopped    = errors.New("program already stopped")
	ErrUnknownState      = errors.New("program state unknown")
	ErrProgramNotReady   = errors.New("program not ready")
	ErrNotWorking        = errors.New("program not working")
	ErrLiveTimeLong      = errors.New("live time is too long")
	ErrPollTimeSmall     = errors.New("poll time is too small")
	ErrInvalidSettings   = errors.New("invalid settings format")
	ErrPollGroupNotExist = errors.New("polltime group does not exist")
)

type DataID uint32
type Undo func() error
type QualityState uint8

type Settings any

type Service interface {
	Run() error
	Stop() error
	Close() error
	Reset() error
	State() uint8
}

type Driver interface {
	TagCreate(id DataID, s Settings) (Undo, error)
	TagDelete(DataID) (Undo, error)
	TagSetValue(id DataID, value []byte) error
	Settings() DriverSettings
}

type Formatter interface {
	String() string
	BytesJSON() ([]byte, error)
}

type DriverSettings struct {
	General Formatter
	Tags    map[DataID]Formatter
}

type Datapoint struct {
	Value     []byte
	Timestamp int64
	Quality   QualityState
}

type Storage interface {
	Create(DataID) (Undo, error)
	Read(DataID) (Datapoint, error)
	Update(DataID, Datapoint) (Undo, error)
	Delete(DataID) (Undo, error)

	ValueUpdater
	List() map[DataID]Datapoint
}

type ValueUpdater interface {
	UpdateValue(DataID, []byte) (Undo, error)
	UpdateQuality(DataID, QualityState) (Undo, error)
}

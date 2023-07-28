package models

const (
	Uncertain QualityState = iota
	Good
	Bad
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

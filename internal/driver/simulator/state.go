package driver

const (
	ready uint8 = iota
	running
	stopping
	reset
	closing
)

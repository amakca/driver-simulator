package driver

const (
	ready uint8 = iota
	running
	stopping
	reset
	closing
)

func (d *simulator) State() uint8 {
	return d.state
}

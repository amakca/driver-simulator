package driver

import (
	m "practice/internal/models"
	"time"
)

func (d *simulator) Run() error {
	switch d.State() {
	case running:
		return errAlreadyRunning
	case closing:
		return errAlreadyClosed
	case stopping:
		d.runSignal()
		d.state = running
		return nil
	case reset:
		return errProgramNotReady
	case ready:
		d.runSignal()
		d.state = running
		for pollTime := range d.pollGroup {
			go func(pollTime time.Duration) {
				d.polling(pollTime)
			}(pollTime)
		}
		return nil
	default:
		return errUnknownState
	}
}

func (d *simulator) Stop() error {
	switch d.State() {
	case stopping:
		return errAlreadyStopped
	case closing:
		return errAlreadyClosed
	case ready, reset:
		return errNotWorking
	case running:
		d.stopSignal()
		d.state = stopping
		for id := range d.tagsSettings {
			d.str.UpdateQuality(id, m.Bad)
		}
		return nil
	default:
		return errUnknownState
	}
}

func (d *simulator) Close() error {
	if d.State() == closing {
		return errAlreadyClosed
	}
	if err := d.dumpConfig(); err != nil {
		return err
	}
	d.closeSignal()

	d.state = closing
	d.shutdown()
	return nil
}

func (d *simulator) Reset() error {
	switch d.State() {
	case closing:
		return errAlreadyClosed
	case reset:
		return errNotWorking
	case stopping, running, ready:
		d.closeSignal()
		d.state = reset
		if err := d.init(d.Settings()); err != nil {
			return err
		}
		d.state = ready
		return nil
	default:
		return errUnknownState
	}
}

func (d *simulator) runSignal() {
	if !IsChanClosable(d.close) {
		d.close = make(chan struct{})
	}
	if !IsChanClosable(d.stop) {
		d.stop = make(chan struct{})
	}
	if IsChanClosable(d.start) {
		close(d.start)
	}
}

func (d *simulator) stopSignal() {
	if !IsChanClosable(d.close) {
		d.close = make(chan struct{})
	}
	if !IsChanClosable(d.start) {
		d.start = make(chan struct{})
	}
	if IsChanClosable(d.stop) {
		close(d.stop)
	}
}

func (d *simulator) closeSignal() {
	if !IsChanClosable(d.start) {
		d.start = make(chan struct{})
	}
	if !IsChanClosable(d.stop) {
		d.stop = make(chan struct{})
	}
	if IsChanClosable(d.close) {
		close(d.close)
	}
}

func (d *simulator) shutdown() {
	if IsChanClosable(d.start) {
		close(d.start)
	}
	if IsChanClosable(d.stop) {
		close(d.stop)
	}
}

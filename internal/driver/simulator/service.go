package driver

import (
	m "practice/internal/models"
	u "practice/internal/utils"
	"time"
)

func (d *simulator) Run() error {
	switch d.State() {
	case running:
		return m.ErrAlreadyRunning
	case closing:
		return m.ErrAlreadyClosed
	case stopping:
		d.runSignal()
		d.state = running
		return nil
	case reset:
		return m.ErrProgramNotReady
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
		return m.ErrUnknownState
	}
}

func (d *simulator) Stop() error {
	switch d.State() {
	case stopping:
		return m.ErrAlreadyStopped
	case closing:
		return m.ErrAlreadyClosed
	case ready, reset:
		return m.ErrNotWorking
	case running:
		d.stopSignal()
		d.state = stopping
		for id := range d.tagsSettings {
			d.storage.UpdateQuality(id, m.BAD)
		}
		return nil
	default:
		return m.ErrUnknownState
	}
}

func (d *simulator) Close() error {
	if d.State() == closing {
		return m.ErrAlreadyClosed
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
		return m.ErrAlreadyClosed
	case reset:
		return m.ErrNotWorking
	case stopping, running, ready:
		d.closeSignal()
		d.state = reset
		if err := d.init(d.Settings()); err != nil {
			return err
		}
		d.state = ready
		return nil
	default:
		return m.ErrUnknownState
	}
}

func (d *simulator) runSignal() {
	if !u.IsChanClosable(d.close) {
		d.close = make(chan struct{})
	}
	if !u.IsChanClosable(d.stop) {
		d.stop = make(chan struct{})
	}
	if u.IsChanClosable(d.start) {
		close(d.start)
	}
}

func (d *simulator) stopSignal() {
	if !u.IsChanClosable(d.close) {
		d.close = make(chan struct{})
	}
	if !u.IsChanClosable(d.start) {
		d.start = make(chan struct{})
	}
	if u.IsChanClosable(d.stop) {
		close(d.stop)
	}
}

func (d *simulator) closeSignal() {
	if !u.IsChanClosable(d.start) {
		d.start = make(chan struct{})
	}
	if !u.IsChanClosable(d.stop) {
		d.stop = make(chan struct{})
	}
	if u.IsChanClosable(d.close) {
		close(d.close)
	}
}

func (d *simulator) shutdown() {
	if u.IsChanClosable(d.start) {
		close(d.start)
	}
	if u.IsChanClosable(d.stop) {
		close(d.stop)
	}
}

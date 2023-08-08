package driver

import (
	"log"
	m "practice/internal/models"
	u "practice/internal/utils"
	"time"
)

func (d *simulator) Run() error {
	switch d.State() {
	case m.RUNNING:
		return m.ErrAlreadyRunning
	case m.CLOSED:
		return m.ErrProgramClosed
	case m.STOPPED:
		d.runSignal()
		d.state = m.RUNNING
		return nil
	case m.READY:
		d.runSignal()
		d.state = m.RUNNING
		for pollTime := range d.pollGroup {
			go func(pollTime time.Duration) {
				if err := d.polling(pollTime); err != nil {
					log.Print(err)
				}
			}(pollTime)
		}
		return nil
	default:
		return m.ErrUnknownState
	}
}

func (d *simulator) Stop() error {
	switch d.State() {
	case m.STOPPED:
		return m.ErrAlreadyStopped
	case m.CLOSED:
		return m.ErrProgramClosed
	case m.READY:
		return m.ErrNotWorking
	case m.RUNNING:
		d.stopSignal()
		d.state = m.STOPPED

		for id := range d.tagsSettings {
			if undo, err := d.storage.UpdateQuality(id, m.QUALITY_BAD); err != nil {
				if undo != nil {
					if err = undo(); err != nil {
						log.Print(err)
					}
					return err
				}
				return err
			}
		}

		return nil
	default:
		return m.ErrUnknownState
	}
}

func (d *simulator) Close() error {
	if d.State() == m.CLOSED {
		return m.ErrProgramClosed
	}
	if err := d.dumpConfig(); err != nil {
		return err
	}
	d.closeSignal()

	d.state = m.CLOSED
	d.shutdown()
	return nil
}

func (d *simulator) Reset() error {
	switch d.State() {
	case m.STOPPED, m.RUNNING, m.READY, m.CLOSED:
		d.closeSignal()
		if err := d.init(d.Settings()); err != nil {
			return err
		}
		d.state = m.READY
		return nil
	default:
		return m.ErrUnknownState
	}
}

func (d *simulator) runSignal() {
	if !u.IsChanClosable(d.stop) {
		d.stop = make(chan struct{})
	}
	if u.IsChanClosable(d.start) {
		close(d.start)
	}
}

func (d *simulator) stopSignal() {
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

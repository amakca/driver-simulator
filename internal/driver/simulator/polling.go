package driver

import (
	g "practice/internal/driver/generator"
	m "practice/internal/models"
	"time"
)

func (d *simulator) controlLiveTime() error {
	go func() {
		delay := time.NewTimer(d.generalSettings.ProgramLiveTime)

		select {
		case <-delay.C:
			delay.Stop()
			d.Close()
		case <-d.close:
			if !delay.Stop() {
				delay.Stop()
			}
		}
	}()
	return nil
}

func (d *simulator) addGenerator(id m.DataID, tagConfig TagSettings) error {
	gen, err := d.genManager.New(tagConfig.GenConfig, d.generalSettings.GenOptimization)
	if err != nil {
		return err
	}

	d.rwmu.Lock()
	defer d.rwmu.Unlock()

	if _, ok := d.pollGroup[tagConfig.PollTime]; !ok {
		d.pollGroup[tagConfig.PollTime] = make(map[m.DataID]*g.Generator)
	}
	d.pollGroup[tagConfig.PollTime][id] = gen

	return nil
}

func (d *simulator) polling(pollTime time.Duration) {

	ticker := time.NewTicker(pollTime)
	for {
		select {
		case <-ticker.C:
			if err := d.updateValue(pollTime); err != nil {
				if err == m.ErrPollGroupNotExist {
					ticker.Stop()
				}
				return
			}
		case <-d.stop:
			ticker.Stop()
			select {
			case <-d.start:
				ticker = time.NewTicker(pollTime)
			case <-d.close:
				d.caseStopGen(pollTime)
				return
			}
		case <-d.close:
			d.caseStopGen(pollTime)
			return
		}

	}
}

func (d *simulator) caseStopGen(pollTime time.Duration) {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()

	for _, gen := range d.pollGroup[pollTime] {
		gen.Stop()
	}

}

func (d *simulator) caseStartGen(pollTime time.Duration) error {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()

	for _, gen := range d.pollGroup[pollTime] {
		if err := gen.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (d *simulator) updateValue(pollTime time.Duration) error {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()

	if _, ok := d.pollGroup[pollTime]; !ok {
		return m.ErrPollGroupNotExist
	}

	for id, gen := range d.pollGroup[pollTime] {
		val := gen.ValueBytes()
		_, err := d.storage.UpdateValue(m.DataID(id), val)
		if err != nil {
			return err
		}
	}

	return nil
}

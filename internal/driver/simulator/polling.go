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
	if _, ok := d.pollGroup[tagConfig.PollTime]; !ok {
		d.pollGroup[tagConfig.PollTime] = make(map[m.DataID]*g.Generator)
	}
	d.pollGroup[tagConfig.PollTime][id] = gen
	d.rwmu.Unlock()

	return nil
}

func (d *simulator) polling(pollTime time.Duration) error {
	d.rwmu.RLock()
	for _, gen := range d.pollGroup[pollTime] {
		if err := gen.Start(); err != nil {
			d.rwmu.RUnlock()
			return err
		}
	}
	d.rwmu.RUnlock()

	ticker := time.NewTicker(pollTime)
	for {
		select {
		case <-ticker.C:
			if err := d.updateValue(pollTime); err != nil {
				if err == m.ErrPollGroupNotExist {
					ticker.Stop()
				}
				return nil
			}
		case <-d.stop:
			ticker.Stop()
			select {
			case <-d.start:
				ticker = time.NewTicker(pollTime)
			case <-d.close:
				d.caseStopGen(pollTime)
				return nil
			}
		case <-d.close:
			d.caseStopGen(pollTime)
			return nil
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

func (d *simulator) updateValue(pollTime time.Duration) error {
	d.rwmu.RLock()

	if _, ok := d.pollGroup[pollTime]; !ok {
		d.rwmu.RUnlock()
		return m.ErrPollGroupNotExist
	}

	for id, gen := range d.pollGroup[pollTime] {
		val := gen.ValueBytes()
		_, err := d.storage.UpdateValue(m.DataID(id), val)
		if err != nil {
			d.rwmu.RUnlock()
			return err
		}
	}

	d.rwmu.RUnlock()
	return nil
}

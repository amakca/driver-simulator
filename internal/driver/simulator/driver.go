package driver

import (
	"fmt"
	g "practice/internal/driver/generator"
	m "practice/internal/models"

	"github.com/pkg/errors"
)

func (d *simulator) TagCreate(id m.DataID, s m.Settings) (m.Undo, error) {
	d.rwmu.Lock()
	defer d.rwmu.Unlock()

	if _, ok := d.tagsSettings[id]; ok {
		return nil, errors.Wrap(m.ErrDataExists,
			fmt.Sprint("id = ", id),
		)
	}

	settings, err := parseTags(s)
	if err != nil {
		return nil, err
	}

	undo := func() error {
		d.rwmu.Lock()
		defer d.rwmu.Unlock()

		return d.deleteTag(id)
	}

	if err := d.createTag(id, settings); err != nil {
		return nil, err
	}

	return undo, nil
}

func (d *simulator) TagDelete(id m.DataID) (m.Undo, error) {
	d.rwmu.Lock()
	defer d.rwmu.Unlock()

	deletedSettings, ok := d.tagsSettings[id]
	if !ok {
		return nil, errors.Wrap(m.ErrDataNotFound,
			fmt.Sprint("id = ", id),
		)
	}

	undo := func() error {
		d.rwmu.Lock()
		defer d.rwmu.Unlock()

		return d.createTag(id, deletedSettings)
	}

	if err := d.deleteTag(id); err != nil {
		if err == g.ErrGenAlreadyStopped {
			return undo, nil
		}
		return nil, err
	}

	return undo, nil
}

func (d *simulator) TagSetValue(id m.DataID, value []byte) error {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()

	for _, intrnl := range d.pollGroup {
		if gen, ok := intrnl[id]; ok {
			err := gen.SetValueBytes(value)
			return err
		}
	}
	return errors.Wrap(m.ErrDataNotFound,
		fmt.Sprint("id = ", id),
	)

}

func (d *simulator) Settings() m.DriverSettings {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()

	settings := m.DriverSettings{
		General: &d.generalSettings,
		Tags:    map[m.DataID]m.Formatter{},
	}

	for k := range d.tagsSettings {
		v := d.tagsSettings[k]
		settings.Tags[k] = &v
	}

	return settings
}

func (d *simulator) createTag(id m.DataID, s TagSettings) error {
	gen, err := d.genManager.New(s.GenConfig, d.generalSettings.GenOptimization)
	if err != nil {
		return err
	}
	if err := gen.Start(); err != nil {
		return err
	}

	d.tagsSettings[id] = s

	_, ok := d.pollGroup[s.PollTime]
	if !ok {
		d.pollGroup[s.PollTime] = make(map[m.DataID]*g.Generator)
	}
	d.pollGroup[s.PollTime][id] = gen

	if !ok {
		go d.polling(s.PollTime)
	}
	return nil
}

func (d *simulator) deleteTag(id m.DataID) error {
	for pollTime, intrnl := range d.pollGroup {
		if gen, ok := intrnl[id]; ok {
			delete(d.tagsSettings, id)
			delete(intrnl, id)

			if len(intrnl) == 0 {
				delete(d.pollGroup, pollTime)
			}

			return gen.Stop()
		}
	}

	return errors.Wrap(m.ErrDataNotFound,
		fmt.Sprint("id = ", id),
	)
}

func (d *simulator) State() uint8 {
	return d.state
}

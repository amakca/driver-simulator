package driver

import (
	g "practice/internal/driver/generator"
	m "practice/internal/models"
	"time"
)

func (d *simulator) TagCreate(id m.DataID, s m.Settings) (m.Undo, error) {
	if _, ok := d.tagsSettings[id]; ok {
		return nil, errDataExist
	}

	settings, err := parseTags(s)
	if err != nil {
		return nil, err
	}

	undo := func() error {
		if err := d.deleteTag(id); err != nil {
			return err
		}
		return nil
	}

	if err := d.createTag(id, settings); err != nil {
		return m.Undo(undo), err
	}

	return m.Undo(undo), nil
}

func (d *simulator) TagDelete(id m.DataID) (m.Undo, error) {
	if _, ok := d.tagsSettings[id]; !ok {
		return nil, errDataNotFound
	}

	oldData := d.tagsSettings[id]

	undo := func() error {
		if err := d.createTag(id, oldData); err != nil {
			return err
		}
		return nil
	}

	if err := d.deleteTag(id); err != nil {
		return m.Undo(undo), err
	}

	return m.Undo(undo), nil
}

func (d *simulator) TagSetValue(id m.DataID, value []byte) error {
	undo, err := d.str.UpdateValue(id, value)
	if err != nil {
		undo()
		return err
	}
	return nil
}

func (d *simulator) Settings() m.DriverSettings {
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
	d.tagsSettings[id] = s

	gen, err := d.genManager.New(s.Settings, d.generalSettings.UseGenManager)
	if err != nil {
		return err
	}

	d.rwmu.Lock()
	_, ok := d.pollGroup[s.PollTime]
	if !ok {
		d.pollGroup[s.PollTime] = make(map[m.DataID]*g.Generator)
	}
	d.pollGroup[s.PollTime][id] = gen
	d.rwmu.Unlock()

	gen.Start()

	if !ok {
		go func(pollTime time.Duration) {
			d.polling(pollTime)
		}(s.PollTime)
	}
	return nil
}

func (d *simulator) deleteTag(id m.DataID) error {
	for pollTime, intrnl := range d.pollGroup {
		if gen, ok := intrnl[id]; ok {
			gen.Stop()
			d.rwmu.Lock()
			delete(intrnl, id)
			d.rwmu.Unlock()
		}

		if len(intrnl) == 0 {
			d.rwmu.Lock()
			delete(d.pollGroup, pollTime)
			d.rwmu.Unlock()
		}

		delete(d.tagsSettings, id)
		return nil
	}
	return errDataNotFound
}

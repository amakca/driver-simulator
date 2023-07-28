package driver

import (
	"encoding/json"
	"os"
	g "practice/internal/driver/generator"
	m "practice/internal/models"
	"sync"
	"time"
)

// Добавить в стору в ундо проверки
// type listGenerators map[m.DataID]*g.Generator
// type pollGroup map[time.Duration]listGenerators

type simulator struct {
	generalSettings GeneralSettings          // Общие настройки драйвера
	tagsSettings    map[m.DataID]TagSettings // Настройки тегов
	// listGenerators  listGenerators           // Лист соотвествия тег-генератор
	// pollGroup       pollGroup                // Группы тегов по времени опроса

	pollGroup map[time.Duration]map[m.DataID]*g.Generator

	genManager *g.GenManager  // Менеджер генераторов
	str        m.ValueUpdater // Хранилище

	start, stop, close chan struct{}

	rwmu  sync.RWMutex
	state uint8
}

func (d *simulator) parseDriverSettings(set m.DriverSettings) (err error) {
	if d.generalSettings, err = parseGeneral(set.General); err != nil {
		return err
	}

	for key, val := range set.Tags {
		if d.tagsSettings[key], err = parseTags(val); err != nil {
			return err
		}
	}

	return nil
}

func New(set m.DriverSettings, str m.ValueUpdater) (*simulator, error) {
	genManager := g.CreateGenManager()

	s := &simulator{
		generalSettings: GeneralSettings{},
		tagsSettings:    make(map[m.DataID]TagSettings),
		pollGroup:       make(map[time.Duration]map[m.DataID]*g.Generator),
		genManager:      genManager,
		str:             str,
		start:           make(chan struct{}),
		state:           ready,
	}

	if err := s.init(set); err != nil {
		return nil, err
	}

	return s, nil
}

func (d *simulator) controlLiveTime() error {
	go func() {
		delay := time.NewTimer(d.generalSettings.MaxLiveTime)

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

func (d *simulator) init(set m.DriverSettings) error {
	if err := d.parseDriverSettings(set); err != nil {
		return err
	}
	for id, settings := range d.tagsSettings {
		gen, err := d.genManager.New(settings.Settings, d.generalSettings.UseGenManager)
		if err != nil {
			return err
		}

		d.rwmu.Lock()
		if _, ok := d.pollGroup[settings.PollTime]; !ok {
			d.pollGroup[settings.PollTime] = make(map[m.DataID]*g.Generator)
		}
		d.pollGroup[settings.PollTime][id] = gen
		d.rwmu.Unlock()
	}

	if err := d.controlLiveTime(); err != nil {
		return err
	}

	return nil
}

func (d *simulator) polling(pollTime time.Duration) error {
	for _, gen := range d.pollGroup[pollTime] {
		if err := gen.Start(); err != nil {
			return err
		}
	}

	ticker := time.NewTicker(pollTime)
	for {
		select {
		case <-ticker.C:
			d.rwmu.RLock()

			if _, ok := d.pollGroup[pollTime]; !ok {
				ticker.Stop()
				d.rwmu.RUnlock()
				return nil
			}

			for id, gen := range d.pollGroup[pollTime] {
				val := gen.ValueBytes()
				_, err := d.str.UpdateValue(m.DataID(id), val)
				if err != nil {
					d.rwmu.RUnlock()
					return err
				}
			}

			d.rwmu.RUnlock()
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

func (d *simulator) dumpConfig() error {
	settings := d.Settings()

	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(settings); err != nil {
		return err
	}
	return nil
}

func (d *simulator) caseStopGen(pollTime time.Duration) {
	d.rwmu.RLock()
	defer d.rwmu.RUnlock()

	for _, gen := range d.pollGroup[pollTime] {
		gen.Stop()
	}

}

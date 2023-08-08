package driver

import (
	"encoding/json"
	"os"
	g "practice/internal/driver/generator"
	m "practice/internal/models"
	u "practice/internal/utils"
	"sync"
	"time"
)

type simulator struct {
	generalSettings GeneralSettings          // Общие настройки драйвера
	tagsSettings    map[m.DataID]TagSettings // Настройки тегов

	pollGroup map[time.Duration]map[m.DataID]*g.Generator

	genManager *g.Manager     // Менеджер генераторов
	storage    m.ValueUpdater // Хранилище

	start, stop, close chan struct{}

	rwmu  sync.RWMutex
	state uint8
}

func New(settings m.DriverSettings, storage m.ValueUpdater) (*simulator, error) {
	genManager, err := g.CreateManager()
	if err != nil {
		return nil, err
	}

	s := &simulator{
		generalSettings: GeneralSettings{},
		tagsSettings:    make(map[m.DataID]TagSettings),
		pollGroup:       make(map[time.Duration]map[m.DataID]*g.Generator),
		genManager:      genManager,
		storage:         storage,
		start:           make(chan struct{}),
	}

	if err := s.init(settings); err != nil {
		return nil, err
	}

	return s, nil
}

func (d *simulator) init(set m.DriverSettings) error {
	if err := d.parseDriverSettings(set); err != nil {
		return err
	}

	for id, tagConfig := range d.tagsSettings {
		if err := d.addGenerator(id, tagConfig); err != nil {
			return err
		}
	}

	if err := d.controlLiveTime(); err != nil {
		return err
	}

	if !u.IsChanClosable(d.close) {
		d.close = make(chan struct{})
	}

	d.state = m.READY
	return nil
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

func (d *simulator) dumpConfig() error {
	settings := d.Settings()

	file, err := os.Create(m.CONFIG_FILE)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(settings)
	return err
}

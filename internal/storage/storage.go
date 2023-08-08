package storage

import (
	"fmt"
	m "practice/internal/models"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type TagStorage struct {
	rwmu sync.RWMutex
	data map[m.DataID]*m.Datapoint
}

func New() (*TagStorage, error) {
	return &TagStorage{
		data: make(map[m.DataID]*m.Datapoint),
	}, nil
}

func (s *TagStorage) Create(id m.DataID) (m.Undo, error) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	if _, ok := s.data[id]; ok {
		return nil, errors.Wrap(m.ErrDataExists,
			fmt.Sprint("id = ", id),
		)
	}

	undo := func() error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		delete(s.data, id)
		return nil
	}

	s.data[id] = &m.Datapoint{
		Quality: m.QUALITY_UNCERTAIN,
	}

	return undo, nil
}

func (s *TagStorage) Read(id m.DataID) (m.Datapoint, error) {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	datapoint, ok := s.data[id]
	if !ok {
		return m.Datapoint{}, errors.Wrap(m.ErrDataNotFound,
			fmt.Sprint("id = ", id),
		)
	}

	return *datapoint, nil
}

func (s *TagStorage) Update(id m.DataID, datapoint m.Datapoint) (m.Undo, error) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	oldDatapoint, ok := s.data[id]
	if !ok {
		return nil, errors.Wrap(m.ErrDataNotFound,
			fmt.Sprint("id = ", id),
		)
	}

	undo := func() error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		s.data[id] = oldDatapoint
		return nil
	}

	s.data[id] = &datapoint

	return undo, nil
}

func (s *TagStorage) Delete(id m.DataID) (m.Undo, error) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	oldDatapoint, ok := s.data[id]
	if !ok {
		return nil, errors.Wrap(m.ErrDataNotFound,
			fmt.Sprint("id = ", id),
		)
	}

	undo := func() error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		s.data[id] = oldDatapoint
		return nil
	}

	delete(s.data, id)

	return undo, nil
}

func (s *TagStorage) List() map[m.DataID]m.Datapoint {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	dataCopy := make(map[m.DataID]m.Datapoint)
	for k, v := range s.data {
		dataCopy[k] = *v
	}

	return dataCopy
}

func (s *TagStorage) UpdateValue(id m.DataID, value []byte) (m.Undo, error) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	datapoint, ok := s.data[id]
	if !ok {
		return nil, errors.Wrap(m.ErrDataNotFound,
			fmt.Sprint("id = ", id),
		)
	}
	oldDatapoint := *datapoint

	undo := func() error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		s.data[id] = &oldDatapoint
		return nil
	}

	s.data[id].Value = value
	s.data[id].Timestamp = time.Now().Unix()

	return undo, nil
}

func (s *TagStorage) UpdateQuality(id m.DataID, state m.QualityState) (m.Undo, error) {
	datapoint, ok := s.data[id]
	if !ok {
		return nil, errors.Wrap(m.ErrDataNotFound,
			fmt.Sprint("id = ", id),
		)
	}
	oldQuality := datapoint.Quality

	undo := func() error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		s.data[id].Quality = oldQuality
		return nil
	}

	s.data[id].Quality = state

	return undo, nil
}

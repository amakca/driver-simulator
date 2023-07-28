package storage

import (
	"errors"
	m "practice/internal/models"
	"sync"
	"time"
)

var (
	errDataExist    = errors.New("data with id already exist")
	errDataNotFound = errors.New("data with id not found")
)

type TagStorage struct {
	rwmu sync.RWMutex
	data map[m.DataID]m.Datapoint
}

func New() (*TagStorage, error) {
	return &TagStorage{
		data: make(map[m.DataID]m.Datapoint),
	}, nil
}

func (s *TagStorage) Create(id m.DataID) (m.Undo, error) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	if _, ok := s.data[id]; ok {
		return nil, errDataExist
	}

	undo := func() error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		if _, ok := s.data[id]; !ok {
			return errDataNotFound
		}

		delete(s.data, id)
		return nil
	}

	s.data[id] = m.Datapoint{
		Quality: m.Uncertain,
	}

	return undo, nil
}

func (s *TagStorage) Read(id m.DataID) (m.Datapoint, error) {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	datapoint, ok := s.data[id]
	if !ok {
		return m.Datapoint{}, errDataNotFound
	}

	return datapoint, nil
}

func (s *TagStorage) Update(id m.DataID, datapoint m.Datapoint) (m.Undo, error) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	oldDatapoint, ok := s.data[id]
	if !ok {
		return nil, errDataNotFound
	}

	undo := func() error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		s.data[id] = oldDatapoint
		return nil
	}

	s.data[id] = datapoint

	return undo, nil
}

func (s *TagStorage) Delete(id m.DataID) (m.Undo, error) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	oldDatapoint, ok := s.data[id]
	if !ok {
		return nil, errDataNotFound
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

	return s.data
}

func (s *TagStorage) UpdateValue(id m.DataID, value []byte) (m.Undo, error) {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	oldDatapoint, ok := s.data[id]
	if !ok {
		return nil, errDataNotFound
	}

	undo := func() error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		s.data[id] = oldDatapoint
		return nil
	}

	s.data[id] = m.Datapoint{
		Value:     value,
		Timestamp: time.Now().Unix(),
		Quality:   m.Good,
	}

	return undo, nil
}

func (s *TagStorage) UpdateQuality(id m.DataID, state m.QualityState) (m.Undo, error) {
	datapoint, ok := s.data[id]
	if !ok {
		return nil, errDataNotFound
	}

	oldQuality := datapoint.Quality

	undo := func() error {
		s.rwmu.Lock()
		defer s.rwmu.Unlock()

		datapoint.Quality = oldQuality
		s.data[id] = datapoint
		return nil
	}

	datapoint.Quality = state
	s.data[id] = datapoint

	return undo, nil
}

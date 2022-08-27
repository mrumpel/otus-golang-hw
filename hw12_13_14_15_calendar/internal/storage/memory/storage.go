package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mrumpel/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	sync.RWMutex
	m map[uuid.UUID]storage.Event
}

func (s *Storage) Connect(ctx context.Context, connectionString string) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, e storage.Event) (uuid.UUID, error) {
	id := uuid.New()

	if e.ID != uuid.Nil {
		s.RLock()
		defer s.RUnlock()
		if _, ok := s.m[e.ID]; ok {
			return uuid.Nil, storage.ErrEventIDAlreadyExist
		}
		id = e.ID
	}

	_, err := s.GetEventList(context.Background(), e.DateStart, e.DateEnd)
	if !errors.Is(err, storage.ErrEventListEmpty) {
		return uuid.Nil, storage.ErrTimeBusy
	}

	s.Lock()
	e.ID = id
	s.m[id] = e
	s.Unlock()

	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id uuid.UUID, e storage.Event) error {
	s.RLock()

	if _, ok := s.m[id]; !ok {
		s.RUnlock()
		return storage.ErrEventNotExist
	}

	if l, err := s.GetEventList(context.Background(), e.DateStart, e.DateEnd); err == nil && (*l)[0].ID != id {
		s.RUnlock()
		return storage.ErrTimeBusy
	}

	s.RUnlock()
	s.Lock()
	e.ID = id
	s.m[id] = e
	s.Unlock()
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	s.RLock()
	if _, ok := s.m[id]; !ok {
		s.RUnlock()
		return storage.ErrEventNotExist
	}
	s.RUnlock()

	s.Lock()
	delete(s.m, id)
	s.Unlock()
	return nil
}

func (s *Storage) GetEventList(ctx context.Context, start time.Time, end time.Time) (*[]storage.Event, error) {
	if start.After(end) {
		start, end = end, start
	}

	res := make([]storage.Event, 0)

	s.RLock()
	for _, e := range s.m {
		if !((e.DateStart.Before(start.Add(1)) && e.DateEnd.Before(start.Add(1))) ||
			(e.DateStart.After(end.Add(-1)) && e.DateEnd.After(end.Add(-1)))) {
			res = append(res, e)
		}
	}
	s.RUnlock()

	if len(res) == 0 {
		return nil, storage.ErrEventListEmpty
	}
	return &res, nil
}

func (s *Storage) GetEventByID(ctx context.Context, id uuid.UUID) (*storage.Event, error) {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.m[id]; !ok {
		return nil, storage.ErrEventNotExist
	}
	res := s.m[id]
	return &res, nil
}

func New() *Storage {
	m := make(map[uuid.UUID]storage.Event)
	return &Storage{
		m: m,
	}
}

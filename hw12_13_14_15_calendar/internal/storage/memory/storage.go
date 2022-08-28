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
	mu     sync.RWMutex
	events map[uuid.UUID]storage.Event
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
		s.mu.RLock()
		defer s.mu.RUnlock()
		if _, ok := s.events[e.ID]; ok {
			return uuid.Nil, storage.ErrEventIDAlreadyExist
		}
		id = e.ID
	}

	_, err := s.GetEventList(context.Background(), e.DateStart, e.DateEnd)
	if !errors.Is(err, storage.ErrEventListEmpty) {
		return uuid.Nil, storage.ErrTimeBusy
	}

	s.mu.Lock()
	e.ID = id
	s.events[id] = e
	s.mu.Unlock()

	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id uuid.UUID, e storage.Event) error {
	s.mu.RLock()

	if _, ok := s.events[id]; !ok {
		s.mu.RUnlock()
		return storage.ErrEventNotExist
	}

	if l, err := s.GetEventList(context.Background(), e.DateStart, e.DateEnd); err == nil && (*l)[0].ID != id {
		s.mu.RUnlock()
		return storage.ErrTimeBusy
	}

	s.mu.RUnlock()
	s.mu.Lock()
	e.ID = id
	s.events[id] = e
	s.mu.Unlock()
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	s.mu.RLock()
	if _, ok := s.events[id]; !ok {
		s.mu.RUnlock()
		return storage.ErrEventNotExist
	}
	s.mu.RUnlock()

	s.mu.Lock()
	delete(s.events, id)
	s.mu.Unlock()
	return nil
}

func (s *Storage) GetEventList(ctx context.Context, start time.Time, end time.Time) (*[]storage.Event, error) {
	if start.After(end) {
		start, end = end, start
	}

	res := make([]storage.Event, 0)

	s.mu.RLock()
	for _, e := range s.events {
		if !((e.DateStart.Before(start.Add(1)) && e.DateEnd.Before(start.Add(1))) ||
			(e.DateStart.After(end.Add(-1)) && e.DateEnd.After(end.Add(-1)))) {
			res = append(res, e)
		}
	}
	s.mu.RUnlock()

	if len(res) == 0 {
		return nil, storage.ErrEventListEmpty
	}
	return &res, nil
}

func (s *Storage) GetEventByID(ctx context.Context, id uuid.UUID) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.events[id]; !ok {
		return nil, storage.ErrEventNotExist
	}
	res := s.events[id]
	return &res, nil
}

func New() *Storage {
	m := make(map[uuid.UUID]storage.Event)
	return &Storage{
		events: m,
	}
}

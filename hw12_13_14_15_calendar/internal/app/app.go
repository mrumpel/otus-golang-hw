package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mrumpel/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage storage.Storage
}

type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
}

type Application interface {
	CreateEvent(context.Context, Event) (uuid.UUID, error)
	UpdateEvent(context.Context, uuid.UUID, Event) error
	DeleteEvent(context.Context, uuid.UUID) error
	GetEventListForDay(time.Time) (*[]Event, error)
	GetEventListForWeek(time.Time) (*[]Event, error)
	GetEventListForMonth(time.Time) (*[]Event, error)
	GetEventList(time.Time, time.Time) (*[]Event, error)
}

type Event struct {
	ID          uuid.UUID
	Title       string
	DateStart   time.Time
	DateEnd     time.Time
	Description string
	OwnerID     string
	AlarmTime   time.Duration
}

func New(logger Logger, storage storage.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, e Event) (uuid.UUID, error) {
	return a.storage.CreateEvent(ctx, storage.Event(e))
}

func (a *App) UpdateEvent(ctx context.Context, id uuid.UUID, e Event) error {
	return a.storage.UpdateEvent(ctx, id, storage.Event(e))
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEventListForDay(t time.Time) (*[]Event, error) {
	return a.GetEventList(t, t.AddDate(0, 0, 1))
}

func (a *App) GetEventListForWeek(t time.Time) (*[]Event, error) {
	return a.GetEventList(t, t.AddDate(0, 0, 7))
}

func (a *App) GetEventListForMonth(t time.Time) (*[]Event, error) {
	return a.GetEventList(t, t.AddDate(0, 1, 0))
}

func (a *App) GetEventList(start, end time.Time) (*[]Event, error) {
	events, err := a.storage.GetEventList(context.Background(), start, end)
	if err != nil {
		return nil, err
	}

	res := make([]Event, 0, len(*events))
	for _, e := range *events {
		res = append(res, Event(e))
	}
	return &res, nil
}

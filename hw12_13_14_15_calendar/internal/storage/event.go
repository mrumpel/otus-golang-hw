package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID     `db:"id"`
	Title       string        `db:"title"`
	DateStart   time.Time     `db:"date_start"`
	DateEnd     time.Time     `db:"date_end"`
	Description string        `db:"description"`
	OwnerID     string        `db:"owner_id"`
	AlarmTime   time.Duration `db:"alarm_time"`
}

type Storage interface {
	Connect(context.Context, string) error
	Close(context.Context) error
	CreateEvent(context.Context, Event) (uuid.UUID, error)
	UpdateEvent(context.Context, uuid.UUID, Event) error
	DeleteEvent(context.Context, uuid.UUID) error
	GetEventList(context.Context, time.Time, time.Time) (*[]Event, error)
	GetEventByID(context.Context, uuid.UUID) (*Event, error)
}

var (
	ErrEventNotExist       = errors.New("event not exist")
	ErrTimeBusy            = errors.New("time is busy")
	ErrEventListEmpty      = errors.New("there is no events")
	ErrEventIDAlreadyExist = errors.New("event already exists")
)

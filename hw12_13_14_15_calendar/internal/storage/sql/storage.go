package sqlstorage

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib" // driver
	"github.com/jmoiron/sqlx"
	"github.com/mrumpel/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	goose "github.com/pressly/goose/v3"
)

type Storage struct {
	db *sqlx.DB
}

var errDBStr = "postgres database error: %w"

//go:embed migrations/*.sql
var embedMigrations embed.FS

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, connectionString string) error {
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return fmt.Errorf(errDBStr, err)
	}
	s.db = db

	err = s.db.Ping()
	if err != nil {
		return fmt.Errorf(errDBStr, err)
	}

	// migration part
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(s.db.DB, "migrations"); err != nil {
		panic(err)
	}

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, e storage.Event) (uuid.UUID, error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}

	_, err := s.GetEventByID(ctx, e.ID)
	if !errors.Is(err, storage.ErrEventNotExist) {
		return uuid.Nil, storage.ErrEventIDAlreadyExist
	}

	_, err = s.GetEventList(ctx, e.DateStart, e.DateEnd)
	if !errors.Is(err, storage.ErrEventListEmpty) {
		return uuid.Nil, storage.ErrTimeBusy
	}

	q := `INSERT INTO public.event (id, title, date_start, date_end, description, owner_id, alarm_time)
			VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err = s.db.ExecContext(ctx, q, e.ID, e.Title, e.DateStart, e.DateEnd, e.Description, e.OwnerID, e.AlarmTime)
	if err != nil {
		return uuid.Nil, fmt.Errorf(errDBStr, err)
	}

	return e.ID, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id uuid.UUID, e storage.Event) error {
	_, err := s.GetEventByID(ctx, e.ID)
	if errors.Is(err, storage.ErrEventNotExist) {
		return storage.ErrEventNotExist
	}

	_, err = s.GetEventList(ctx, e.DateStart, e.DateEnd)
	if !errors.Is(err, storage.ErrEventListEmpty) {
		return storage.ErrTimeBusy
	}

	q := `UPDATE public.event 
	SET
    	title = $2,
		date_start = $3, 
		date_end = $4,
		description = $5,
		owner_id = $6,
		alarm_time = $7
	WHERE event.id = $1`
	_, err = s.db.ExecContext(ctx, q, e.ID, e.Title, e.DateStart, e.DateEnd, e.Description, e.OwnerID, e.AlarmTime)
	if err != nil {
		return fmt.Errorf(errDBStr, err)
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	_, err := s.GetEventByID(ctx, id)
	if errors.Is(err, storage.ErrEventNotExist) {
		return storage.ErrEventNotExist
	}

	q := `DELETE FROM public.event WHERE event.id = $1`
	_, err = s.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf(errDBStr, err)
	}

	return nil
}

func (s *Storage) GetEventList(ctx context.Context, start time.Time, end time.Time) (*[]storage.Event, error) {
	if start.After(end) {
		start, end = end, start
	}

	res := make([]storage.Event, 0)

	q := `SELECT * FROM public.event 
	WHERE NOT (
		(date_start < $1 AND date_end < $1) 
		OR (date_start > $2 AND date_end > $2))`
	err := s.db.SelectContext(ctx, &res, q, start, end)
	if err != nil {
		return nil, fmt.Errorf(errDBStr, err)
	}
	if len(res) == 0 {
		return nil, storage.ErrEventListEmpty
	}

	return &res, nil
}

func (s *Storage) GetEventByID(ctx context.Context, id uuid.UUID) (*storage.Event, error) {
	res := make([]storage.Event, 0)

	q := `SELECT * FROM public.event WHERE id = $1`
	err := s.db.SelectContext(ctx, &res, q, id)
	if err != nil {
		return nil, fmt.Errorf(errDBStr, err)
	}
	if len(res) == 0 {
		return nil, storage.ErrEventNotExist
	}

	return &res[0], nil
}

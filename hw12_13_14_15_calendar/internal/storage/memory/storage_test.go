package memorystorage

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mrumpel/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("storage creating", func(t *testing.T) {
		s := *New()
		require.NotNil(t, s.events)
	})

	t.Run("create and check event", func(t *testing.T) {
		s := New()

		id, err := s.CreateEvent(context.Background(), storage.Event{
			Title:     "Hello World Event",
			DateStart: time.Now(),
			DateEnd:   time.Now().Add(1 * time.Hour),
			OwnerID:   "007",
		})
		require.NoError(t, err)

		e, err := s.GetEventByID(context.Background(), id)
		require.NoError(t, err)
		require.NotEmpty(t, e)
		require.Contains(t, e.Title, "Hello")
	})

	t.Run("time busy", func(t *testing.T) {
		s := New()
		t1 := time.Now()
		t2 := t1.Add(1 * time.Hour)
		_, err := s.CreateEvent(context.Background(), storage.Event{
			Title:     "Hello World Event",
			DateStart: t1,
			DateEnd:   t2,
		})
		require.NoError(t, err)

		id2, err := s.CreateEvent(context.Background(), storage.Event{
			Title:     "Excess event",
			DateStart: t1.Add(15 * time.Minute),
			DateEnd:   t1.Add(15 * time.Minute),
		})
		require.True(t, errors.Is(err, storage.ErrTimeBusy))
		require.Equal(t, id2, uuid.Nil)
	})

	t.Run("update and delete event", func(t *testing.T) {
		s := New()
		t1 := time.Now()
		t2 := t1.Add(1 * time.Hour)
		id, err := s.CreateEvent(context.Background(), storage.Event{
			Title:     "Hello World Event",
			DateStart: t1,
			DateEnd:   t2,
		})
		require.NoError(t, err)

		// update part
		err = s.UpdateEvent(context.Background(), id, storage.Event{
			Title:     "Updated event",
			DateStart: t1.Add(15 * time.Minute),
			DateEnd:   t2.Add(15 * time.Minute),
		})
		require.NoError(t, err)

		e, err := s.GetEventByID(context.Background(), id)
		require.NoError(t, err)
		require.Contains(t, e.Title, "Updated")
		require.True(t, e.DateStart.After(t1))
		require.True(t, e.DateEnd.After(t2))

		// delete part
		err = s.DeleteEvent(context.Background(), id)
		require.NoError(t, err)

		_, err = s.GetEventByID(context.Background(), id)
		require.True(t, errors.Is(err, storage.ErrEventNotExist))
	})

	t.Run("concurrency test", func(t *testing.T) {
		s := New()
		t1 := time.Now()
		t2 := t1.Add(1 * time.Hour)
		var wg sync.WaitGroup
		wg.Add(1000)
		for i := 0; i < 1000; i++ {
			t1 = t1.Add(1 * time.Hour)
			t2 = t2.Add(1 * time.Hour)
			go func(t1, t2 time.Time) {
				defer wg.Done()
				_, err := s.CreateEvent(context.Background(), storage.Event{
					Title:     "Hello World Event",
					DateStart: t1,
					DateEnd:   t2,
				})
				require.NoError(t, err)
			}(t1, t2)
		}
		wg.Wait()
	})
}

func TestStorage_GetEventList(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		start, end time.Time
		wantErr    bool
		count      int
	}{
		{
			name:    "no events",
			start:   now.AddDate(-1, 0, 0),
			end:     now.Add(-1),
			wantErr: true,
		},
		{
			name:    "whole time",
			start:   now,
			end:     now.Add(4 * time.Hour),
			wantErr: false,
			count:   3,
		},
		{
			name:    "events crossing",
			start:   now.Add(30 * time.Minute),
			end:     now.Add(150 * time.Minute),
			wantErr: false,
			count:   2,
		},
		{
			name:    "one moment at the border",
			start:   now.Add(3 * time.Hour),
			end:     now.Add(3 * time.Hour),
			wantErr: true,
			count:   0,
		},
		{
			name:    "one moment in the event",
			start:   now.Add(30 * time.Minute),
			end:     now.Add(30 * time.Minute),
			wantErr: false,
			count:   1,
		},
	}

	events := []struct {
		t1, t2 time.Time
	}{
		{now, now.Add(1 * time.Hour)},
		{now.Add(2 * time.Hour), now.Add(3 * time.Hour)},
		{now.Add(3 * time.Hour), now.Add(4 * time.Hour)},
	}

	s := New()

	for i := range events {
		_, err := s.CreateEvent(context.Background(), storage.Event{
			Title:     "event #" + strconv.Itoa(i),
			DateStart: events[i].t1,
			DateEnd:   events[i].t2,
		})
		require.NoError(t, err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetEventList(context.Background(), tt.start, tt.end)
			if err != nil {
				require.True(t, errors.Is(err, storage.ErrEventListEmpty) && tt.wantErr)
				return
			}
			require.Equal(t, tt.count, len(*got))
		})
	}
}

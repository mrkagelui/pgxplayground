package player

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
)

//func TestStore_Info(t *testing.T) {
//	tests := []struct {
//		name     string
//		bomb     string
//		playerID int64
//		want     PlayerInfo
//		check    errCheckFunc
//	}{
//		{
//			name:     "no such player",
//			playerID: 1,
//			want:     PlayerInfo{},
//			check:    errorSays("getting player: no rows in result set"),
//		},
//		{
//			name:     "food error",
//			bomb:     "ALTER TABLE IF EXISTS food_changes RENAME COLUMN change TO delta",
//			playerID: 3,
//			want:     PlayerInfo{},
//			check:    errorSays(`getting player: scany: rows final error: ERROR: relation "players" does not exist (SQLSTATE 42P01)`),
//		},
//		{
//			name:     "food annihilated",
//			bomb:     "DROP TABLE IF EXISTS food_changes",
//			playerID: 3,
//			want:     PlayerInfo{},
//			check:    errorSays(`getting food: scany: rows final error: ERROR: relation "food_changes" does not exist (SQLSTATE 42P01)`),
//		},
//		{
//			name:     "all good",
//			playerID: 3,
//			want: PlayerInfo{
//				ID:    3,
//				Age:   1,
//				Food:  150,
//				Wood:  250,
//				Gold:  550,
//				Stone: 240,
//			},
//			check: noError,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Store{pool}
//			tx, err := pool.Begin(context.Background())
//			noError(t, err)
//			s.db = tx
//
//			if tt.bomb != "" {
//				_, e := tx.Exec(context.Background(), tt.bomb)
//				noError(t, e)
//			}
//
//			got, err := s.Info(context.Background(), tt.playerID)
//			tt.check(t, err)
//			equals(t, got, tt.want)
//
//			noError(t, tx.Rollback(context.Background()))
//		})
//	}
//}

type errCheckFunc func(*testing.T, error)

func noError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func errorSays(wantMsg string) errCheckFunc {
	return func(t *testing.T, err error) {
		if err.Error() != wantMsg {
			t.Errorf("got err %v, want %v", err.Error(), wantMsg)
		}
	}
}

func errorIs(wantErr error) errCheckFunc {
	return func(t *testing.T, err error) {
		if !errors.Is(err, wantErr) {
			t.Errorf("%v is not in %v", wantErr, err)
		}
	}
}

func errorIsSays(wantErr error, wantMsg string) errCheckFunc {
	return func(t *testing.T, err error) {
		errorIs(wantErr)(t, err)
		errorSays(wantMsg)(t, err)
	}
}

func equals[T comparable](t *testing.T, got, want T) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestStore_Get(t *testing.T) {
	tests := []struct {
		name     string
		playerID int64
		want     DBPlayer
		check    errCheckFunc
	}{
		{
			name:     "no such player",
			playerID: 0,
			check:    errorSays("collecting: no rows in result set"),
		},
		{
			name:     "all good",
			playerID: 3,
			want: DBPlayer{
				ID:        3,
				CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
				UpdatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
				Age:       1,
			},
			check: noError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				db: pool,
			}
			got, err := s.Get(context.Background(), tt.playerID)
			tt.check(t, err)
			equals(t, got, tt.want)
		})
	}
}

type fetcher interface {
	fetch(context.Context) ([]DBPlayer, error)
}

type withScany struct{}

func (w withScany) fetch(ctx context.Context) ([]DBPlayer, error) {
	var r []DBPlayer
	return r, pgxscan.Select(ctx, pool, &r, `SELECT id, created_at, updated_at, age FROM players`)
}

func (w withScany) Get(ctx context.Context, id int64) (DBPlayer, error) {
	var r DBPlayer
	return r, pgxscan.Get(ctx, pool, &r, `SELECT id, created_at, updated_at, age FROM players WHERE id = $1`, id)
}

func BenchmarkFetch(b *testing.B) {
	for _, tt := range []struct {
		name string
		f    fetcher
	}{
		{
			name: "pgx",
			f:    &Store{pool},
		},
		{
			name: "with scany",
			f:    withScany{},
		},
	} {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := tt.f.fetch(context.Background())
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

type getter interface {
	Get(ctx context.Context, id int64) (DBPlayer, error)
}

func BenchmarkGet(b *testing.B) {
	for _, tt := range []struct {
		name string
		f    getter
	}{
		{
			name: "pgx",
			f:    &Store{pool},
		},
		{
			name: "with scany",
			f:    withScany{},
		},
	} {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := tt.f.Get(context.Background(), 1)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

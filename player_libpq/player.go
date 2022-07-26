// Package player provides information of the player
package player

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

// Store contains everything needed to provide player info
type Store struct {
	db sqlx.ExtContext
}

// PlayerInfo is what the name implies
type PlayerInfo struct {
	ID    int64
	Age   int
	Food  int
	Wood  int
	Gold  int
	Stone int
}

// Info returns PlayerInfo given the ID
func (s *Store) Info(ctx context.Context, playerID int64) (PlayerInfo, error) {
	type dbPlayer struct {
		ID  int64 `db:"id"`
		Age int   `db:"age"`
	}

	var p dbPlayer
	if err := sqlx.GetContext(ctx, s.db, &p, `SELECT id, age FROM players WHERE id = $1`, playerID); err != nil {
		return PlayerInfo{}, fmt.Errorf("getting player: %w", err)
	}

	var food, wood, gold, stone int
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := sqlx.GetContext(gCtx, s.db, &food, `SELECT SUM(change) FROM food_changes WHERE player_id = $1`, playerID); err != nil {
			return fmt.Errorf("getting food: %v", err)
		}
		return nil
	})
	g.Go(func() error {
		if err := sqlx.GetContext(gCtx, s.db, &wood, `SELECT SUM(change) FROM wood_changes WHERE player_id = $1`, playerID); err != nil {
			return fmt.Errorf("getting wood: %v", err)
		}
		return nil
	})
	g.Go(func() error {
		if err := sqlx.GetContext(gCtx, s.db, &gold, `SELECT SUM(change) FROM gold_changes WHERE player_id = $1`, playerID); err != nil {
			return fmt.Errorf("getting gold: %v", err)
		}
		return nil
	})
	g.Go(func() error {
		if err := sqlx.GetContext(gCtx, s.db, &stone, `SELECT SUM(change) FROM stone_changes WHERE player_id = $1`, playerID); err != nil {
			return fmt.Errorf("getting stone: %v", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return PlayerInfo{}, err
	}

	return PlayerInfo{
		ID:    p.ID,
		Age:   p.Age,
		Food:  food,
		Wood:  wood,
		Gold:  gold,
		Stone: stone,
	}, nil
}

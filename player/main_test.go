package player

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
)

var (
	//go:embed seeds/seed.sql
	seedSQL string
	// shared db connection
	pool *pgxpool.Pool
)

func TestMain(m *testing.M) {
	var err error
	pool, err = getSeededDB(seedSQL)
	if err != nil {
		log.Println("getting db:", err)
		os.Exit(1)
	}
	defer pool.Close()

	m.Run()
}

type dbConfig struct {
	PgHost     string `default:"localhost"`
	PgPort     int    `default:"5442"`
	PgUser     string `default:"user"`
	PgPassword string `default:"password"`
	PgDatabase string `default:"pgxplayground"`
}

// getSeededDB returns a db connection which is seeded
func getSeededDB(seedSQL string) (*pgxpool.Pool, error) {
	var dbConf dbConfig
	if err := envconfig.Process("TEST", &dbConf); err != nil {
		return nil, fmt.Errorf("reading env: %v", err)
	}

	p, err := open(dbConf)
	if err != nil {
		return nil, fmt.Errorf("opening db: %v", err)
	}

	if err := seed(p, seedSQL); err != nil {
		return nil, fmt.Errorf("seeding: %v", err)
	}

	return p, nil
}

func open(cfg dbConfig) (*pgxpool.Pool, error) {
	sslMode := "disable"

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.PgUser, cfg.PgPassword),
		Host:     fmt.Sprintf("%v:%v", cfg.PgHost, cfg.PgPort),
		Path:     cfg.PgDatabase,
		RawQuery: q.Encode(),
	}

	p, err := pgxpool.New(context.Background(), u.String())
	if err != nil {
		return nil, fmt.Errorf("connect: %v", err)
	}

	return p, nil
}

func seed(db *pgxpool.Pool, seedSQL string) error {
	if seedSQL == "" {
		return nil
	}
	if _, err := db.Exec(context.Background(), seedSQL); err != nil {
		return err
	}
	return nil
}

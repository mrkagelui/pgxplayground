package player

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

var (
	//go:embed seeds/seed.sql
	seedSQL string
	// shared db connection
	dbConn *sqlx.DB
)

func TestMain(m *testing.M) {
	var err error
	dbConn, err = getSeededDB(seedSQL)
	if err != nil {
		log.Println("getting db:", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	os.Exit(m.Run())
}

type dbConfig struct {
	PgHost     string `default:"localhost"`
	PgPort     int    `default:"5442"`
	PgUser     string `default:"user"`
	PgPassword string `default:"password"`
	PgDatabase string `default:"pgxplayground"`
}

// getSeededDB returns a db connection which is seeded
func getSeededDB(seedSQL string) (*sqlx.DB, error) {
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

func open(cfg dbConfig) (*sqlx.DB, error) {
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

	db, err := sqlx.Connect("postgres", u.String())
	if err != nil {
		return nil, fmt.Errorf("connect: %v", err)
	}

	return db, nil
}

func seed(db *sqlx.DB, seedSQL string) error {
	if seedSQL == "" {
		return nil
	}
	if _, err := db.ExecContext(context.Background(), seedSQL); err != nil {
		return err
	}
	return nil
}

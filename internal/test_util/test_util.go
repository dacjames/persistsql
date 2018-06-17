package test_util

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
)

func WithTestDB(t *testing.T, block func(*sql.DB)) {
	pool, err := dockertest.NewPool("")
	require.Nil(t, err)

	pool.MaxWait = 5 * time.Second

	resource, err := pool.Run("postgres", "10.4", []string{"POSTGRES_PASSWORD=postgres"})
	require.Nil(t, err)

	defer func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Could not purge resource: %s", err)
		}
	}()

	var db *sql.DB
	if err := pool.Retry(func() error {
		var err error
		conn := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/postgres?sslmode=disable", resource.GetPort("5432/tcp"))

		db, err = sql.Open("postgres", conn)
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	block(db)
}

func Migrated(t *testing.T, block func(*sql.DB)) func(*sql.DB) {
	return func(db *sql.DB) {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		require.Nil(t, err)

		m, err := migrate.NewWithDatabaseInstance("file://../../migrations", "postgres", driver)
		require.Nil(t, err)

		err = m.Up()
		require.Nil(t, err)

		block(db)

		err = m.Down()
		require.Nil(t, err)
	}
}

func WithMigratedDB(t *testing.T, block func(*sql.DB)) {
	WithTestDB(t, Migrated(t, block))
}

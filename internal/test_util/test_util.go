package test_util

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

type dockerContext struct {
	name     string
	resource *dockertest.Resource
	pool     *dockertest.Pool
}

type Teardowner interface {
	Teardown()
}

func (t dockerContext) Teardown() {
	if t.resource == nil {
		return
	}
	if err := t.pool.Purge(t.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err.Error())
	}
}

func SetupTestDB() (*sql.DB, Teardowner) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not setup dockertest pool: %s", err.Error())
	}

	pool.MaxWait = 5 * time.Second

	_, filename, _, _ := runtime.Caller(1)
	key := strings.Split(path.Base(filename), ".")[0]
	name := fmt.Sprintf("%s_%s", "persistsql", key)

	existing, err := pool.Client.ListContainers(docker.ListContainersOptions{
		Filters: map[string][]string{
			"name": []string{name},
		},
	})
	if err != nil {
		log.Fatalf("Could not list existing containers: %s", err.Error())
	}

	var resource *dockertest.Resource
	var port string
	if len(existing) == 0 {
		resource, err = pool.RunWithOptions(&dockertest.RunOptions{
			Labels: map[string]string{
				"persistsql": "test",
				"test":       key,
			},
			Name:       name,
			Repository: "postgres",
			Tag:        "10.4",
			Env:        []string{"POSTGRES_PASSWORD=postgres"},
		})
		if err != nil {
			log.Fatalf("Could not run postgres container: %s", err.Error())
		}
		port = resource.GetPort("5432/tcp")
	} else {
		c, err := pool.Client.InspectContainer(existing[0].ID)
		if err != nil {
			log.Fatalf("Could not inspect existing postgres container: %s", err.Error())
		}

		if p, ok := c.NetworkSettings.Ports[docker.Port("5432/tcp")]; ok && len(p) > 0 {
			port = p[0].HostPort
		} else {
			log.Fatalf("Could locate host port for 5432/tcp on container %s", existing[0].ID)
		}

	}

	var db *sql.DB
	if err := pool.Retry(func() error {
		var err error
		conn := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/postgres?sslmode=disable", port)

		db, err = sql.Open("postgres", conn)
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	return db, dockerContext{
		name:     name,
		pool:     pool,
		resource: resource,
	}
}

func CleanupTestContainers() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not setup dockertest pool: %s", err.Error())
	}

	existing, err := pool.Client.ListContainers(docker.ListContainersOptions{
		Filters: map[string][]string{
			"label": []string{"persistsql=test"},
		},
	})

	for _, e := range existing {
		err := pool.Client.RemoveContainer(docker.RemoveContainerOptions{ID: e.ID, Force: true, RemoveVolumes: true})
		if err != nil {
			log.Printf("Failed to remove container %s: %s", e.ID, err.Error())
		}
	}
}

type migrateContext struct {
	m *migrate.Migrate
}

func (c migrateContext) Teardown() {
	if err := c.m.Down(); err != nil {
		log.Fatalf("Could not down migrate database: %s", err.Error())
	}
}

func Migrate(db *sql.DB) Teardowner {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not setup postgres migration instance: %s", err.Error())
	}

	// Safe because GOPATH is always set when running tests
	gopath := os.Getenv("GOPATH")
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/src/github.com/dacjames/persistsql/migrations", gopath), "postgres", driver)
	if err != nil {
		log.Fatalf("Failed to load database migrations: %s", err.Error())
	}

	err = m.Up()
	if err != nil {
		log.Fatalf("Failed to apply database migrations: %s", err.Error())
	}

	return migrateContext{m: m}
}

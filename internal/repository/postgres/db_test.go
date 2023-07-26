package postgres

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"log"
	"os"
	"testing"
	"time"
)

var testRepository FilmsRepository

func TestMain(m *testing.M) {

	testRepository = FilmsRepository{TransactionManager: NewTransactionManager()}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=test",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://postgres:postgres@%s/test?sslmode=disable", hostAndPort)

	resource.Expire(360)

	if err = pool.Retry(func() error {
		testRepository.Pool, err = pgxpool.New(context.Background(), databaseUrl)
		if err != nil {
			return err
		}
		return testRepository.Pool.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	pool.MaxWait = 30 * time.Second

	PerformMigrationUp(databaseUrl)

	code := m.Run()

	os.Exit(code)

}

func PerformMigrationUp(dbURL string) {
	m, err := migrate.New("file://../migrations", dbURL)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if err = m.Up(); err != nil {
		log.Fatalf(err.Error())
	}
}

package reindexer

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models/api"
	dbModel "github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/models/general"
	"github.com/google/uuid"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/restream/reindexer/v3"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

type Status int

const (
	SUCCESS Status = iota
	FAILURE
)

var testRepository CacheRepository

func TestMain(m *testing.M) {

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	namespace := "test"
	testRepository = CacheRepository{namespace: namespace}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "reindexer/reindexer",
		Env:        []string{},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("6534/tcp")

	dsn := fmt.Sprintf("cproto://%s/%s", hostAndPort, namespace)
	log.Println(dsn)
	resource.Expire(360)

	if err = pool.Retry(func() error {
		testRepository.db = reindexer.NewReindex(dsn, reindexer.WithCreateDBIfMissing())
		return testRepository.db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = testRepository.db.OpenNamespace(namespace, reindexer.DefaultNamespaceOptions(), dbModel.Film{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	pool.MaxWait = 30 * time.Second

	code := m.Run()

	os.Exit(code)

}

func TestCacheRepository_Add(t *testing.T) {
	film := dbModel.NewFilm(uuid.New(),
		fmt.Sprintf("%s%d", t.Name(), rand.Int()),
		&general.Crew{
			Actors:    []*general.Actor{},
			Producers: []*general.Producer{},
		},
		2010,
		10.99,
		"comedy",
	)

	_, err := testRepository.Add(context.Background(), *film)
	assert.Equal(t, nil, err)
}

func TestCacheRepository_FindByName(t *testing.T) {
	nameToBeFound := fmt.Sprintf("%s%d", t.Name(), rand.Int())
	film := dbModel.NewFilm(uuid.New(),
		nameToBeFound,
		&general.Crew{
			Actors: []*general.Actor{},
			Producers: []*general.Producer{
				{Person: general.Person{Name: nameToBeFound, Birthdate: "1970-07-30", Gender: "m", ID: uuid.New()}},
			},
		},
		2023,
		400000,
		"biography",
	)
	_, err := testRepository.Add(context.Background(), *film)
	assert.Equal(t, nil, err)

	testTable := []struct {
		name                  string
		lengthOfReceivedSlice int
	}{
		{
			name:                  nameToBeFound,
			lengthOfReceivedSlice: 1,
		}, {
			name:                  "",
			lengthOfReceivedSlice: 0,
		},
	}

	for _, test := range testTable {
		films, err := testRepository.FindByName(context.Background(), test.name)
		assert.Equal(t, nil, err)
		assert.Equal(t, test.lengthOfReceivedSlice, len(films))
	}
}

func TestCacheRepository_Delete(t *testing.T) {
	nameToBeVerified := fmt.Sprintf("%s%d", t.Name(), rand.Int())
	film := dbModel.NewFilm(uuid.New(),
		nameToBeVerified,
		&general.Crew{
			Actors: []*general.Actor{},
			Producers: []*general.Producer{
				{Person: general.Person{Name: nameToBeVerified, Birthdate: "1970-07-30", Gender: "m", ID: uuid.New()}},
			},
		},
		2004,
		1000,
		"biography",
	)
	_, err := testRepository.Add(context.Background(), *film)
	assert.Equal(t, nil, err)

	testTable := []struct {
		request api.DeleteRequest
		status  Status
	}{
		{
			request: api.DeleteRequest{Title: film.Title,
				Genre:        film.Genre,
				ReleasedYear: film.ReleasedYear},
			status: SUCCESS,
		}, {
			request: api.DeleteRequest{},
			status:  FAILURE,
		},
	}

	for _, test := range testTable {
		err = testRepository.Delete(context.Background(), test.request)
		if test.status == SUCCESS {
			assert.Equal(t, nil, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}

func TestCacheRepository_DeleteCachedWithCtx(t *testing.T) {
	nameToBeVerified := fmt.Sprintf("%s%d", t.Name(), rand.Int())
	film := dbModel.NewFilm(uuid.New(),
		nameToBeVerified,
		&general.Crew{
			Actors: []*general.Actor{},
			Producers: []*general.Producer{
				{Person: general.Person{Name: nameToBeVerified, Birthdate: "1970-07-30", Gender: "m", ID: uuid.New()}},
			},
		},
		2004,
		1000,
		"biography",
	)
	_, err := testRepository.Add(context.Background(), *film)
	assert.Equal(t, nil, err)

	testTable := []struct {
		request api.DeleteRequest
		status  Status
	}{
		{
			request: api.DeleteRequest{Title: film.Title,
				Genre:        film.Genre,
				ReleasedYear: film.ReleasedYear},
			status: SUCCESS,
		}, {
			request: api.DeleteRequest{},
			status:  FAILURE,
		},
	}

	for _, test := range testTable {
		ctxWithValues := context.WithValue(context.Background(), "request", test.request)
		err = testRepository.DeleteCachedWithCtx(ctxWithValues)
		if test.status == SUCCESS {
			assert.Equal(t, nil, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}

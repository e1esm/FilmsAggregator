package postgres

import (
	"context"
	"fmt"
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/models/general"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestFilmsRepository_FindAll(t *testing.T) {

	testTable := []struct {
		film   db.Film
		status Status
	}{
		{
			film: *db.NewFilm(uuid.New(),
				fmt.Sprintf("%v%d", t.Name(), rand.Int()), &general.Crew{Producers: []*general.Producer{}, Actors: []*general.Actor{}},
				2000,
				199999.99,
				"comedy",
			),
			status: SUCCESS,
		},
		{
			film: *db.NewFilm(uuid.New(),
				fmt.Sprintf("%v%d", t.Name(), rand.Int()), &general.Crew{Producers: []*general.Producer{}, Actors: []*general.Actor{}},
				2000,
				199999.99,
				"comedy",
			),
		},
		{
			film: *db.NewFilm(uuid.New(),
				fmt.Sprintf("%v%d", t.Name(), rand.Int()), &general.Crew{Producers: []*general.Producer{}, Actors: []*general.Actor{}},
				1000,
				999.99,
				"detective",
			),
			status: FAIL,
		},
	}
	toBeFound := 0
	for _, test := range testTable {
		if test.status == SUCCESS {
			_, err := testRepository.Add(context.Background(), test.film)
			assert.Equal(t, nil, err)
			toBeFound++
		} else {
			continue
		}
	}

	films, err := testRepository.FindAll(context.Background())
	assert.Equal(t, nil, err)
	i := 0
	actuallyFound := 0

	for _, film := range films {
		if film.Title == testTable[i].film.Title {
			actuallyFound++
		}
		i++
	}

	assert.Equal(t, toBeFound, actuallyFound)
}

func TestFilmsRepository_FindByName(t *testing.T) {
	film := db.NewFilm(uuid.New(),
		fmt.Sprintf("%v%d", t.Name(), rand.Int()), &general.Crew{Producers: []*general.Producer{}, Actors: []*general.Actor{}},
		1000,
		999.99,
		"detective",
	)

	_, err := testRepository.Add(context.Background(), *film)
	assert.Equal(t, nil, err)
	testTable := []struct {
		title  string
		status Status
		err    error
	}{
		{
			title:  film.Title,
			status: SUCCESS,
		}, {
			title:  "",
			status: FAIL,
			err:    pgx.ErrNoRows,
		},
	}

	for _, test := range testTable {
		_, receivedErr := testRepository.FindByName(context.Background(), test.title)
		assert.Equal(t, test.err, receivedErr)
	}
}

func TestFilmsRepository_Verify(t *testing.T) {
	film := db.NewFilm(uuid.New(),
		fmt.Sprintf("%v%d", t.Name(), rand.Int()), &general.Crew{Producers: []*general.Producer{}, Actors: []*general.Actor{}},
		1000,
		999.99,
		"history",
	)

	_, err := testRepository.Add(context.Background(), *film)
	assert.Equal(t, nil, err)
	testTable := []struct {
		film   db.Film
		result bool
	}{
		{
			film: *db.NewFilm(uuid.New(),
				film.Title, &general.Crew{Producers: []*general.Producer{}, Actors: []*general.Actor{}},
				1000,
				999.99,
				"history",
			),
			result: true,
		},
		{
			film: *db.NewFilm(uuid.New(),
				fmt.Sprintf("%v%d", t.Name(), rand.Int()), &general.Crew{Producers: []*general.Producer{}, Actors: []*general.Actor{}},
				1000,
				999.99,
				"history",
			),
			result: false,
		}, {
			film: *db.NewFilm(uuid.New(),
				"", &general.Crew{Producers: []*general.Producer{}, Actors: []*general.Actor{}},
				1997,
				20000000,
				"romance",
			),
			result: false,
		},
	}

	for _, test := range testTable {
		doesSatisfy := testRepository.Verify(context.Background(), &test.film)
		assert.Equal(t, test.result, doesSatisfy)
	}
}

func TestFilmsRepository_FindFilmsByActor(t *testing.T) {
	film := db.NewFilm(uuid.New(),
		"Barbie",
		&general.Crew{
			Actors: []*general.Actor{
				{
					Person: general.Person{Name: "Margot Robbie", ID: uuid.New(), Birthdate: "1990-03-24", Gender: "f"},
				},
			},
			Producers: []*general.Producer{},
		},
		2023,
		399999.99,
		"comedy",
	)

	_, err := testRepository.Add(context.Background(), *film)
	assert.Equal(t, nil, err)

	testTable := []struct {
		name string
		err  error
	}{
		{
			name: "Margot Robbie",
			err:  nil,
		},
		{
			name: "",
			err:  pgx.ErrNoRows,
		},
	}

	for _, test := range testTable {
		_, err := testRepository.FindFilmsByActor(context.Background(), test.name)
		assert.Equal(t, test.err, err)
	}
}

func TestFilmsRepository_FindFilmsByProducer(t *testing.T) {
	film := db.NewFilm(uuid.New(),
		"Oppenheimer",
		&general.Crew{
			Actors: []*general.Actor{},
			Producers: []*general.Producer{
				{Person: general.Person{Name: "Christopher Nolan", Birthdate: "1970-07-30", Gender: "m", ID: uuid.New()}},
			},
		},
		2023,
		400000,
		"biography",
	)
	_, err := testRepository.Add(context.Background(), *film)
	assert.Equal(t, nil, err)

	testTable := []struct {
		name string
		err  error
	}{
		{
			name: "Christopher Nolan",
			err:  nil,
		},
		{
			name: "",
			err:  pgx.ErrNoRows,
		},
	}

	for _, test := range testTable {
		_, err = testRepository.FindFilmsByProducer(context.Background(), test.name)
		assert.Equal(t, test.err, err)
	}
}

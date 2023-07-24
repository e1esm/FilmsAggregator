package uuid

import (
	"github.com/e1esm/FilmsAggregator/internal/models/db"
	"github.com/e1esm/FilmsAggregator/internal/models/general"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Status int

const (
	OK Status = iota
	Fail
)

func TestUUIDGenerator_Generate(t *testing.T) {
	idgenerator := UUIDGenerator{}
	testTable := []struct {
		inputID     uuid.UUID
		generatedID uuid.UUID
	}{
		{
			inputID:     uuid.UUID{},
			generatedID: idgenerator.Generate(),
		},
		{
			inputID:     idgenerator.Generate(),
			generatedID: idgenerator.Generate(),
		},
	}

	for _, test := range testTable {
		assert.NotEqual(t, test.inputID, test.generatedID)
	}
}

func TestUUIDGenerator_GenerateUUIDs(t *testing.T) {
	idgenerator := UUIDGenerator{}
	testTable := []struct {
		inputFilm db.Film
		status    Status
	}{
		{
			inputFilm: db.Film{Title: "XXX",
				Crew: &general.Crew{Actors: []*general.Actor{&general.Actor{Role: "Y"}},
					Producers: []*general.Producer{&general.Producer{Person: general.Person{Name: "YY"}}}}},
			status: OK,
		}, {
			inputFilm: db.Film{Title: "YYY",
				Crew: nil},
			status: Fail,
		},
	}
	emptyID := uuid.UUID{}

	for _, test := range testTable {
		isOK := true
		receivedFilm := idgenerator.GenerateUUIDs(test.inputFilm)
		if receivedFilm.Crew == nil {
			isOK = false
			if test.status == Fail {
				assert.Equal(t, false, isOK)
			}
			return
		}
		for _, producer := range receivedFilm.Crew.Actors {
			if producer.ID == emptyID {
				isOK = false
			}
		}
		for _, actor := range receivedFilm.Crew.Actors {
			if actor.ID == emptyID {
				isOK = false
			}
		}

		if test.status == OK {
			assert.Equal(t, true, isOK)
		} else {
			assert.Equal(t, false, isOK)
		}
	}

}

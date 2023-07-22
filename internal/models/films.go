package models

import (
	"github.com/google/uuid"
	"time"
)

type Film struct {
	ID           uuid.UUID `json:"-" reindex:"-"`
	CacheID      int64     `json:"-" reindex:"id,,pk"`
	Title        string    `json:"title" reindex:"title,tree"`
	Crew         Crew      `json:"crew"`
	ReleasedYear int       `json:"released_year" reindex:"released_year"`
	Revenue      float64   `json:"revenue" reindex:"revenue"`
	CacheTime    time.Time `json:"-" reindex:"cache_time"`
}

type Crew struct {
	Actors    []Actor    `json:"actors"`
	Producers []Producer `json:"producers"`
}

type Person struct {
	ID        uuid.UUID `json:"-" reindex:"-"`
	Name      string    `json:"name" reindex:"name"`
	Birthdate string    `json:"birthdate" reindex:"birthdate"`
	Gender    string    `json:"gender" reindex:"gender"`
}

func NewPerson(id uuid.UUID, name string, birthdate string, gender string) Person {
	return Person{ID: id, Name: name, Birthdate: birthdate, Gender: gender}
}

type Producer struct {
	Person
}

func NewProducer(id uuid.UUID, name string, birthdate string, gender string) *Producer {
	return &Producer{NewPerson(id, name, birthdate, gender)}
}

type Actor struct {
	Person
	Role string `json:"role" reindex:"role"`
}

func NewActor(id uuid.UUID, name string, birthdate string, gender string, role string) *Actor {
	return &Actor{Role: role, Person: NewPerson(id, name, birthdate, gender)}
}

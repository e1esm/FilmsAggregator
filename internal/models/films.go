package models

import "github.com/google/uuid"

type Film struct {
	ID           uuid.UUID `json:"-"`
	Title        string    `json:"title"`
	Crew         Crew      `json:"crew"`
	ReleasedYear int       `json:"released_year"`
	Revenue      float64   `json:"revenue"`
}

type Crew struct {
	Actors    []Actor    `json:"actors"`
	Producers []Producer `json:"producers"`
}

type Person struct {
	ID        uuid.UUID `json:"-"`
	Name      string    `json:"name"`
	Birthdate string    `json:"birthdate"`
	Gender    string    `json:"gender"`
}

type Producer struct {
	Person
}

type Actor struct {
	Person
	Role string `json:"role"`
}

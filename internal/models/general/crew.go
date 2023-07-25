package general

import "github.com/google/uuid"

// Crew model info
// @Description The model of the crew which took part in shooting the show
type Crew struct {
	Actors    []*Actor    `json:"actors"`    // All actors took a part in the show
	Producers []*Producer `json:"producers"` // All producers that took a part in the show
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

// Producer model info
// @Description producer of the show
type Producer struct {
	Person
}

func NewProducer(id uuid.UUID, name string, birthdate string, gender string) *Producer {
	return &Producer{NewPerson(id, name, birthdate, gender)}
}

// Actor model info
// @Description Actor that took a part in the show
type Actor struct {
	Person
	Role string `json:"role" reindex:"role"` // Actor's role in the show
}

func NewActor(id uuid.UUID, name string, birthdate string, gender string, role string) *Actor {
	return &Actor{Role: role, Person: NewPerson(id, name, birthdate, gender)}
}

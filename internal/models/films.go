package models

type Producer struct {
	Name  string `json:"name"`
	Films []Film `json:"films"`
}

type Film struct {
	Title        string  `json:"title"`
	Actors       []Actor `json:"actors"`
	ReleasedYear int     `json:"released_year"`
	Revenue      float64 `json:"revenue"`
}

type Actor struct {
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
}

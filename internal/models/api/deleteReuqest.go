package api

type DeleteRequest struct {
	Genre        string `json:"genre"`
	Title        string `json:"title"`
	ReleasedYear int    `json:"released_year"`
}

package api

// DeleteRequest model info
// @Description request, according to which a film is going to be deleted from the database.
// @Description It has 3 filters: Genre, Title and ReleasedYear
type DeleteRequest struct {
	Genre        string `json:"genre"`         //A genre of the show to be deleted
	Title        string `json:"title"`         //A title of the show to be deleted
	ReleasedYear int    `json:"released_year"` // A year of release of a show to be deleted
}
